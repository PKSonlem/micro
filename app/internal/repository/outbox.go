package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type OutboxEvent struct {
	ID        int       `db:"id"`
	EventType string    `db:"event_type"`
	FlatID    int       `db:"flat_id"`
	HouseID   int       `db:"house_id"`
	Email     string    `db:"email"`
	Message   string    `db:"message"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *Repository) GetPendingOutboxEvents(ctx context.Context, limit int) ([]OutboxEvent, error) {
	query := r.qb.
		Select("*").
		From(outboxTable).
		Where(sq.Eq{"status": "pending"}).
		OrderBy("created_at ASC").
		Limit(uint64(limit)).
		Suffix("FOR UPDATE SKIP LOCKED") // чтобы воркеры не взяли одно и то же

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	rows, err := r.conn.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying pending events: %w", err)
	}
	defer rows.Close()

	var events []OutboxEvent
	for rows.Next() {
		var event OutboxEvent
		err = rows.Scan(
			&event.ID,
			&event.EventType,
			&event.FlatID,
			&event.HouseID,
			&event.Email,
			&event.Message,
			&event.Status,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return events, nil
}

func (r *Repository) MarkEventsProcessing(ctx context.Context, eventIDs []int) error {
	if len(eventIDs) == 0 {
		return nil
	}

	query := r.qb.
		Update(outboxTable).
		Set("status", "processing").
		Where(sq.Eq{"id": eventIDs})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	_, err = r.conn.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("error updating events status to processing: %w", err)
	}

	return nil
}

func (r *Repository) MarkEventsSent(ctx context.Context, eventIDs []int) error {
	if len(eventIDs) == 0 {
		return nil
	}

	query := r.qb.
		Update(outboxTable).
		Set("status", "sent").
		Where(sq.Eq{"id": eventIDs})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	_, err = r.conn.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("error updating events status to sent: %w", err)
	}

	return nil
}

func (r *Repository) MarkEventsFailed(ctx context.Context, eventIDs []int) error {
	if len(eventIDs) == 0 {
		return nil
	}

	query := r.qb.
		Update(outboxTable).
		Set("status", "failed").
		Where(sq.Eq{"id": eventIDs})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	_, err = r.conn.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("error updating events status to failed: %w", err)
	}

	return nil
}

func (r *Repository) ResetStuckEvents(ctx context.Context, stuckTimeout time.Duration) (int, error) {
	stuckThreshold := time.Now().Add(-stuckTimeout)

	query := r.qb.
		Update(outboxTable).
		Set("status", "pending").
		Where(sq.And{
			sq.Eq{"status": "processing"},
			sq.Lt{"created_at": stuckThreshold},
		})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("error building query: %w", err)
	}

	result, err := r.conn.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return 0, fmt.Errorf("error resetting stuck events: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting rows affected: %w", err)
	}

	return int(rowsAffected), nil
}
