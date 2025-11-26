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

func (u *Usecase) Handle(ctx context.Context, flat entity.Flat) (*entity.Flat, error) {
	return u.repo.CreateFlat(ctx, flat)
}
