package create_house

import (
	"github.com/timurzdev/mentorship-test-task/internal/deps"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

type Usecase struct {
	repo   Repository
	logger deps.Logger
}

func NewUsecase(repo Repository, logger deps.Logger) *Usecase {
	return &Usecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *Usecase) Handle(house entity.House) error {
	return u.repo.CreateHouse(house)
}
