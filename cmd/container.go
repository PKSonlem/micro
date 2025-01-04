package cmd

import (
	"context"
	"log"

	embPg "github.com/fergusstrange/embedded-postgres"
	"github.com/jmoiron/sqlx"
	"github.com/timurzdev/mentorship-test-task/migrations"
)

// контейне внешних зависимостей приложения
// тут мы инициализируем все инфраструктурные зависимости
type Container struct {
	gCtx             context.Context
	configuration    *configuration
	db               *sqlx.DB
	embeddedPostgres *embPg.EmbeddedPostgres
	migrator         *migrations.Migrator
	logger           *Logger
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

func (e *Container) GetEmbeddedPostgres() *embPg.EmbeddedPostgres {
	if e.embeddedPostgres == nil {
		e.embeddedPostgres = embPg.NewDatabase(
			e.configuration.
				GetPostgresConfiguration().
				GetEmbeddedPostgresConfig(),
		)
	}

	return e.embeddedPostgres
}

func (e *Container) GetLogger() *Logger {
	if e.logger == nil {
		e.logger = NewLogger()
	}

	return e.logger
}

func (e *Container) GetMigrator() *migrations.Migrator {
	if e.migrator == nil {
		e.migrator = migrations.NewMigrator(
			e.configuration.
				GetPostgresConfiguration().
				GetMigrateConnectionString(),
		)
	}

	return e.migrator
}
