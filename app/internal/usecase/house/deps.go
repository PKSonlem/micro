package house

import (
	"context"

	"github.com/PKSonlem/micro/internal/entity"
)

//go:generate mockgen -source=deps.go -destination=mock/deps.go -package=mock
type repository interface {
	CreateHouse(ctx context.Context, house entity.House) (*entity.House, error)
	GetHouseFlats(ctx context.Context, houseId int, isModerator bool) ([]entity.Flat, error)
	CreateSubscription(ctx context.Context, houseID int, email string) error
}
