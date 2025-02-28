package create_house

import (
	"context"

	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

//go:generate mockgen -source=deps.go -destination=mock/deps.go -package=mock
type repository interface {
	CreateHouse(ctx context.Context, flat entity.House) (*entity.House, error)
}
