package cmd

import (
	"github.com/jmoiron/sqlx"
	"github.com/timurzdev/mentorship-test-task/internal/repository"
)

// контейнер внутренних зависимостей
type Internal struct {
	repo *repository.Repository
}

func NewInternal() *Internal {
	return &Internal{}
}

func (i *Internal) GetRepository(conn *sqlx.DB) *repository.Repository {
	if i.repo == nil {
		i.repo = repository.NewRepository(conn)
	}

	return i.repo
}
