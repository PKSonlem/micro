package repository

import (
	"context"
	"fmt"
	"time"

	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/helpers"
)

const (
	houseTable = "house"
)

func (r *Repository) CreateHouse(ctx context.Context, house entity.House) error {
	return sqlxTransaction(ctx, r.conn, func(tx *sqlx.Tx) error {
		return r.createHouseTx(ctx, house, tx)
	})
}

func (r *Repository) createHouseTx(ctx context.Context, house entity.House, tx *sqlx.Tx) error {
	insertMap := map[string]interface{}{
		"address":    house.Address,
		"year":       house.Year,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	// проверка, т.к это необязательный параметр по условиям задачи
	if house.Developer != nil {
		insertMap["developer"] = helpers.FromPtr(house.Developer)
	}

	sql, args, err := r.qb.
		Insert(houseTable).
		SetMap(insertMap).
		ToSql()

	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Join(entity.ErrorCreatingHouse, err)
	}

	return nil
}
