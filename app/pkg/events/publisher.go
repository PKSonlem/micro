package events

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type AnalyticsRepository interface {
	InsertBusinessEvent(ctx context.Context, event BusinessEvent) error
}

type RolesProvider interface {
	GetUserID(ctx context.Context) (string, error)
}

type Logger interface {
	Error(ctx context.Context, err error, args ...any)
	Info(ctx context.Context, msg string, args ...any)
}

type EventType string

const (
	EventUserRegistered    EventType = "user_registered"
	EventUserLogin         EventType = "user_login"
	EventHouseCreated      EventType = "house_created"
	EventHouseViewed       EventType = "house_viewed"
	EventHouseSubscription EventType = "house_subscription"
	EventFlatCreated       EventType = "flat_created"
	EventFlatUpdated       EventType = "flat_updated"
)

type BusinessEvent struct {
	EventID    string
	EventType  EventType
	UserID     string
	EntityID   int
	EntityType string
	Metadata   map[string]interface{}
	Status     string
	Error      string
}

type Publisher struct {
	analyticsRepo AnalyticsRepository
	provider      RolesProvider
	logger        Logger
}

func NewPublisher(analyticsRepo AnalyticsRepository, provider RolesProvider, logger Logger) *Publisher {
	return &Publisher{
		analyticsRepo: analyticsRepo,
		provider:      provider,
		logger:        logger,
	}
}

func (p *Publisher) Publish(ctx context.Context, event BusinessEvent) error {
	go func() {
		if err := p.analyticsRepo.InsertBusinessEvent(ctx, event); err != nil {
			p.logger.Error(ctx, fmt.Errorf("failed to publish event %s: %w", event.EventType, err))
		}
	}()

	return nil
}

func (p *Publisher) PublishWithUser(ctx context.Context, eventType EventType, entityID int, entityType string, metadata map[string]interface{}) {
	userID, err := p.provider.GetUserID(ctx)
	if err != nil {
		userID = ""
	}

	event := BusinessEvent{
		EventID:    uuid.New().String(),
		EventType:  eventType,
		UserID:     userID,
		EntityID:   entityID,
		EntityType: entityType,
		Metadata:   metadata,
		Status:     "success",
		Error:      "",
	}

	_ = p.Publish(ctx, event)
}
