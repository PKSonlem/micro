package migrations

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var fs embed.FS

type Migrator struct {
	*migrate.Migrate
}

func NewMigrator(databaseUrl string) *Migrator {
	d, err := iofs.New(fs, ".")
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, databaseUrl)
	if err != nil {
		panic(err)
	}

	return &Migrator{m}
}

func (m *Migrator) MigrateUp() error {
	var err error
	if err = m.Up(); errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}

func (m *Migrator) MigrateDown() error {
	return m.MigrateDown()
}
