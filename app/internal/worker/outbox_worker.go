package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/timurzdev/mentorship-test-task/internal/repository"
	"github.com/timurzdev/mentorship-test-task/pkg/logger"
	"github.com/timurzdev/mentorship-test-task/pkg/sender"
)

type OutboxWorker struct {
	repo          *repository.Repository
	sender        *sender.Sender
	logger        *logger.Logger
	analyticsRepo *repository.AnalyticsRepository
	interval      time.Duration
	limit         int
	stuckTimeout  time.Duration
}

func NewOutboxWorker(repo *repository.Repository, sender *sender.Sender, logger *logger.Logger, analyticsRepo *repository.AnalyticsRepository, interval time.Duration, limit int) *OutboxWorker {
	return &OutboxWorker{
		repo:          repo,
		sender:        sender,
		logger:        logger,
		analyticsRepo: analyticsRepo,
		interval:      interval,
		limit:         limit,
		stuckTimeout:  5 * time.Minute,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	w.logger.Info(ctx, "Outbox worker started")

	w.resetStuckEvents(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.resetStuckEvents(ctx)
			w.processBatch(ctx)

		case <-ctx.Done():
			w.logger.Info(ctx, "Outbox worker shutting down gracefully")
			return
		}
	}
}

func (w *OutboxWorker) resetStuckEvents(ctx context.Context) {
	count, err := w.repo.ResetStuckEvents(ctx, w.stuckTimeout)
	if err != nil {
		w.logger.Error(ctx, fmt.Errorf("error resetting stuck events: %w", err))
		return
	}
	if count > 0 {
		w.logger.Info(ctx, fmt.Sprintf("Reset %d stuck events back to pending", count))
	}
}

func (w *OutboxWorker) processBatch(ctx context.Context) {
	events, err := w.repo.GetPendingOutboxEvents(ctx, w.limit)
	if err != nil {
		w.logger.Error(ctx, fmt.Errorf("error getting pending events: %w", err))
		return
	}

	if len(events) == 0 {
		return
	}

	w.logger.Info(ctx, fmt.Sprintf("Processing %d events", len(events)))

	eventIDs := make([]int, len(events))
	for i, event := range events {
		eventIDs[i] = event.ID
	}

	if err = w.repo.MarkEventsProcessing(ctx, eventIDs); err != nil {
		w.logger.Error(ctx, fmt.Errorf("error marking events as processing: %w", err))
		return
	}

	var successIDs []int
	var failedIDs []int

	for _, event := range events {
		if err = w.sender.SendEmail(ctx, event.Email, event.Message); err != nil {
			w.logger.Error(ctx, fmt.Errorf("failed to send event %d to %s: %w", event.ID, event.Email, err))
			failedIDs = append(failedIDs, event.ID)
		} else {
			w.logger.Info(ctx, fmt.Sprintf("Event %d sent to %s", event.ID, event.Email))
			successIDs = append(successIDs, event.ID)
		}
	}

	if len(successIDs) > 0 {
		if err = w.repo.MarkEventsSent(ctx, successIDs); err != nil {
			w.logger.Error(ctx, fmt.Errorf("error marking events as sent: %w", err))
		}
	}

	if len(failedIDs) > 0 {
		if err = w.repo.MarkEventsFailed(ctx, failedIDs); err != nil {
			w.logger.Error(ctx, fmt.Errorf("error marking events as failed: %w", err))
		}
	}

	w.logger.Info(ctx, fmt.Sprintf("Batch processed: %d successful, %d failed", len(successIDs), len(failedIDs)))

	if len(events) > 0 {
		statuses := make(map[int]string)
		for _, id := range successIDs {
			statuses[id] = "sent"
		}
		for _, id := range failedIDs {
			statuses[id] = "failed"
		}

		err = w.analyticsRepo.InsertOutboxEvents(ctx, events, statuses)
		if err != nil {
			w.logger.Error(ctx, fmt.Errorf("failed to insert events to clickhouse: %w", err))
		}
	}
}
