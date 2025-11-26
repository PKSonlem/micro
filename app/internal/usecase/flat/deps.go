package flat

import (
	"context"

	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

//go:generate mockgen -source=deps.go -destination=mock/deps.go -package=mock
type repository interface {
	CreateFlat(ctx context.Context, flat entity.Flat) (*entity.Flat, error)
}
