package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/timurzdev/mentorship-test-task/pkg/events"
)

type AnalyticsRepository struct {
	conn driver.Conn
}

func NewAnalyticsRepository(conn driver.Conn) *AnalyticsRepository {
	return &AnalyticsRepository{conn: conn}
}

func (a *AnalyticsRepository) InsertOutboxEvents(ctx context.Context, events []OutboxEvent, statuses map[int]string) error {
	batch, err := a.conn.PrepareBatch(ctx, "INSERT INTO outbox_events (event_id, event_type, flat_id, house_id, email, status, created_at)")
	if err != nil {
		return fmt.Errorf("error preparing batching: %w", err)
	}

	for _, event := range events {
		status := statuses[event.ID]
		err = batch.Append(
			int32(event.ID),
			event.EventType,
			int32(event.FlatID),
			int32(event.HouseID),
			event.Email,
			status,
			event.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("error appending to batch: %w", err)
		}
	}

	err = batch.Send()
	if err != nil {
		return fmt.Errorf("error sending batch: %w", err)
	}

	return nil
}

func (a *AnalyticsRepository) InsertBusinessEvent(ctx context.Context, event events.BusinessEvent) error {
	batch, err := a.conn.PrepareBatch(ctx, "INSERT INTO business_events")
	if err != nil {
		return fmt.Errorf("error preparing batch: %w", err)
	}

	metadataJSON, _ := json.Marshal(event.Metadata)

	err = batch.Append(
		event.EventID,
		event.EventType,
		event.UserID,
		int32(event.EntityID),
		event.EntityType,
		string(metadataJSON),
		event.Status,
		event.Error,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("error appending to batch: %w", err)
	}

	err = batch.Send()
	if err != nil {
		return fmt.Errorf("error sending batch: %w", err)
	}

	return nil
}
