package create_house

import (
	"context"

	"github.com/timurzdev/mentorship-test-task/internal/deps"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

type Usecase struct {
	repo   repository
	logger deps.Logger
}

func NewUsecase(repo repository, logger deps.Logger) *Usecase {
	return &Usecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *Usecase) Handle(ctx context.Context, house entity.House) error {
	return u.repo.CreateHouse(ctx, house)
}
