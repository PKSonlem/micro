package flat

import (
	"context"

	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/pkg/events"
)

type Usecase struct {
	repo           repository
	eventPublisher *events.Publisher
}

func NewUsecase(repo repository, eventPublisher *events.Publisher) *Usecase {
	return &Usecase{
		repo:           repo,
		eventPublisher: eventPublisher,
	}
}

func (u *Usecase) HandleCreateFlat(ctx context.Context, flat entity.Flat) (*entity.Flat, error) {
	result, err := u.repo.CreateFlat(ctx, flat)
	if err != nil {
		return nil, err
	}

	u.eventPublisher.PublishWithUser(ctx, events.EventFlatCreated, result.ID, "flat", map[string]interface{}{
		"house_id": result.HouseID,
		"price":    result.Price,
		"rooms":    result.Rooms,
		"status":   result.Status,
	})

	return result, nil
}

func (u *Usecase) HandleUpdateStatus(ctx context.Context, flatID int, status string) (*entity.Flat, error) {
	result, err := u.repo.UpdateModeratorFlat(ctx, flatID, status)
	if err != nil {
		return nil, err
	}

	u.eventPublisher.PublishWithUser(ctx, events.EventFlatUpdated, result.ID, "flat", map[string]interface{}{
		"flat_id":  result.ID,
		"status":   result.Status,
		"house_id": result.HouseID,
	})

	return result, nil
}
