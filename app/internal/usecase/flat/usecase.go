package flat

import (
	"context"

	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

type Usecase struct {
	repo repository
}

func NewUsecase(repo repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) HandleCreateFlat(ctx context.Context, flat entity.Flat) (*entity.Flat, error) {
	return u.repo.CreateFlat(ctx, flat)
}

func (u *Usecase) HandleUpdateStatus(ctx context.Context, flatID int, status string) (*entity.Flat, error) {
	return u.repo.UpdateModeratorFlat(ctx, flatID, status)
}
