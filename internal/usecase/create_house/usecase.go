package create_house

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

func (u *Usecase) Handle(ctx context.Context, house entity.House) (*entity.House, error) {
	return u.repo.CreateHouse(ctx, house)
}
