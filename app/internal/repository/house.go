package repository

import (
	"context"
	"fmt"
	"time"

	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

const (
	houseTable = "house"
)

type houseRow struct {
	Id        int       `db:"id"`
	Address   string    `db:"address"`
	Year      int       `db:"year"`
	Developer *string   `db:"developer"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// действуем по простому правилу - экспортируемый метод - транзакция
func (r *Repository) CreateHouse(ctx context.Context, house entity.House) (*entity.House, error) {
	var res *entity.House
	var err, txErr error

	txErr = sqlxTransaction(ctx, r.conn, func(tx *sqlx.Tx) error {
		res, err = r.createHouseTx(ctx, house, tx)
		return err
	})
	if txErr != nil {
		return nil, txErr
	}

	return res, nil
}

// неэкспортируемый файл - атомраная операция, которую мы можем поместить в любую транзакцию
func (r *Repository) createHouseTx(ctx context.Context, house entity.House, tx *sqlx.Tx) (*entity.House, error) {
	insertMap := map[string]any{
		"address":    house.Address,
		"year":       house.Year,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	// проверка, т.к это необязательный параметр по условиям задачи
	if house.Developer != nil {
		insertMap["developer"] = *house.Developer
	}

	sql, args, err := r.qb.
		Insert(houseTable).
		SetMap(insertMap).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	var row houseRow
	err = tx.GetContext(ctx, &row, sql, args...)
	if err != nil {
		return nil, errors.Join(entity.ErrorCreatingHouse, err)
	}

	result := &entity.House{
		ID:        row.Id,
		Address:   row.Address,
		Year:      row.Year,
		Developer: row.Developer,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	return result, nil
}
