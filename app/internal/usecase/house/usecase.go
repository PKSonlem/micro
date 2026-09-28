package house

import (
	"context"

	"github.com/PKSonlem/micro/internal/entity"
	"github.com/PKSonlem/micro/pkg/events"
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

func (u *Usecase) Handle(ctx context.Context, house entity.House) (*entity.House, error) {
	result, err := u.repo.CreateHouse(ctx, house)
	if err != nil {
		return nil, err
	}

	u.eventPublisher.PublishWithUser(ctx, events.EventHouseCreated, result.ID, "house", map[string]interface{}{
		"address":   result.Address,
		"year":      result.Year,
		"developer": result.Developer,
	})

	return result, nil
}

func (u *Usecase) HandleGetHouseFlatID(ctx context.Context, houseID int, isModerator bool) ([]entity.Flat, error) {
	result, err := u.repo.GetHouseFlats(ctx, houseID, isModerator)
	if err != nil {
		return nil, err
	}

	u.eventPublisher.PublishWithUser(ctx, events.EventHouseViewed, houseID, "house", map[string]interface{}{
		"house_id":     houseID,
		"is_moderator": isModerator,
		"flats_count":  len(result),
	})

	return result, nil

}

func (u *Usecase) Subscribe(ctx context.Context, houseID int, email string) error {
	err := u.repo.CreateSubscription(ctx, houseID, email)
	if err != nil {
		return err
	}

	u.eventPublisher.PublishWithUser(ctx, events.EventHouseSubscription, houseID, "subscribtion", map[string]interface{}{
		"house_id": houseID,
		"email":    email,
	})

	return nil
}
