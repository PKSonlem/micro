package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

const (
	flatTable = "flat"
)

type flatRow struct {
	ID         int       `db:"id"`
	HouseID    int       `db:"house_id"`
	FlatNumber int       `db:"flat_number"`
	Price      int       `db:"price"`
	Rooms      int       `db:"rooms"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// Транзакция
func (r *Repository) CreateFlat(ctx context.Context, flat entity.Flat) (*entity.Flat, error) {
	var res *entity.Flat
	var err, txErr error

	txErr = sqlxTransaction(ctx, r.conn, func(tx *sqlx.Tx) error {
		res, err = r.createFlatTx(ctx, flat, tx)
		return err
	})
	if txErr != nil {
		return nil, txErr
	}

	return res, nil
}

// Атомарная операция для квартиры
func (r *Repository) createFlatTx(ctx context.Context, flat entity.Flat, tx *sqlx.Tx) (*entity.Flat, error) {
	var maxFlatNumber sql.NullInt64
	err := tx.QueryRowContext(ctx,
		"SELECT MAX(flat_number) FROM flat WHERE house_id = $1",
		flat.HouseID).Scan(&maxFlatNumber)

	if err != nil {
		return nil, fmt.Errorf("error getting flat number: %w", err)
	}

	nextFlatNumber := 1
	if maxFlatNumber.Valid {
		nextFlatNumber = int(maxFlatNumber.Int64) + 1
	}

	insertMap := map[string]any{
		"house_id":    flat.HouseID,
		"flat_number": nextFlatNumber,
		"price":       flat.Price,
		"rooms":       flat.Rooms,
		"status":      "created",
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	sql, args, err := r.qb.
		Insert(flatTable).
		SetMap(insertMap).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	updateSQL, updateArgs, err := r.qb. // добавляем время создания квартиры в струкуре дом
						Update(houseTable).
						Set("add_flat", time.Now()).
						Where(sq.Eq{"id": flat.HouseID}). // sq функция для создания учловия равенства в SQl-запросах (WHERE id = <значение flat.HouseID>)
						ToSql()

	if err != nil {
		return nil, fmt.Errorf("error build update query: %w", err)
	}

	var row flatRow
	err = tx.GetContext(ctx, &row, sql, args...)
	if err != nil {
		return nil, errors.Join(entity.ErrorCreatingFlat, err)
	}

	_, err = tx.ExecContext(ctx, updateSQL, updateArgs...)
	if err != nil {
		return nil, fmt.Errorf("error update house add_flat: %w", err)
	}

	result := &entity.Flat{
		ID:         row.ID,
		HouseID:    row.HouseID,
		FlatNumber: row.FlatNumber,
		Price:      row.Price,
		Rooms:      row.Rooms,
		Status:     row.Status,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}

	return result, nil
}

func (r *Repository) UpdateModeratorFlat(ctx context.Context, flatID int, status string) (*entity.Flat, error) {
	var res *entity.Flat
	var err, txErr error

	txErr = sqlxTransaction(ctx, r.conn, func(tx *sqlx.Tx) error {
		res, err = r.updateModeratorFlat(ctx, flatID, status, tx)
		return err
	})
	if txErr != nil {
		return nil, txErr
	}

	return res, nil
}

func (r *Repository) updateModeratorFlat(ctx context.Context, flatID int, status string, tx *sqlx.Tx) (*entity.Flat, error) {
	insertMap := map[string]any{
		"status":     status,
		"updated_at": time.Now(),
	}

	sql, args, err := r.qb.
		Update(flatTable).
		SetMap(insertMap).
		Where(sq.Eq{"id": flatID}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	var row flatRow
	err = tx.GetContext(ctx, &row, sql, args...)
	if err != nil {
		return nil, errors.Join(entity.ErrorUpdateModeratorFlat, err)
	}

	result := &entity.Flat{
		ID:         row.ID,
		HouseID:    row.HouseID,
		FlatNumber: row.FlatNumber,
		Price:      row.Price,
		Rooms:      row.Rooms,
		Status:     row.Status,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}

	return result, nil
}
