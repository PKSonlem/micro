package cmd

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

// контейне внешних зависимостей приложения
// тут мы инициализируем все инфраструктурные зависимости
type Container struct {
	gCtx          context.Context
	configuration *configuration
	db            *sqlx.DB
	logger        *Logger
}

func NewContainer() *Container {
	return &Container{
		configuration: newFromEnv(),
	}
}

// для доступа внутренних зависимостей к конфигурации
func (e *Container) GetConfiguration() *configuration {
	return e.configuration
}

func (e *Container) GetGlobalContext() context.Context {
	if e.gCtx == nil {
		e.gCtx = context.Background()
	}

	return e.gCtx
}

func (e *Container) GetPostgres() *sqlx.DB {
	if e.db == nil {
		var err error
		e.db, err = NewSqlxConn(e.configuration.GetPostgresConfiguration())
		if err != nil {
			log.Fatal(err)
		}
	}

	return e.db
}

func (e *Container) GetLogger() *Logger {
	if e.logger == nil {
		e.logger = NewLogger()
	}

	return e.logger
}
