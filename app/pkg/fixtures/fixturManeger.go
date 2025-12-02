package fixtures

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"

	"github.com/jmoiron/sqlx"
)

//go:embed fixtures/*.sql
var fixtureFS embed.FS

type FixtureManager struct {
	db *sqlx.DB
}

func NewFixtureManager(db *sqlx.DB) *FixtureManager {
	return &FixtureManager{
		db: db,
	}
}

func (fm *FixtureManager) LoadFixtures() error {
	entries, err := fs.ReadDir(fixtureFS, "fixtures")
	if err != nil {
		return fmt.Errorf("error read dir: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}

		filePath := "fixtures/" + ent.Name()

		query, readErr := fs.ReadFile(fixtureFS, filePath)
		if readErr != nil {
			return fmt.Errorf("error read file: %w", readErr)
		}

		sql := string(query)
		_, err = fm.db.ExecContext(context.Background(), sql)
		if err != nil {
			return fmt.Errorf("error sql exec: %w", err)
		}
	}

	return nil
}
