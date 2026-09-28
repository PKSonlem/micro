package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/PKSonlem/micro/internal/entity"
)

const (
	houseTable = "house"
)

type houseRow struct {
	ID        int        `db:"id"`
	Address   string     `db:"address"`
	Year      int        `db:"year"`
	Developer *string    `db:"developer"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	AddFlat   *time.Time `db:"add_flat"`
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

	if house.AddFlat != nil {
		insertMap["add_flat"] = *house.AddFlat
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
		ID:        row.ID,
		Address:   row.Address,
		Year:      row.Year,
		Developer: row.Developer,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		AddFlat:   row.AddFlat,
	}

	return result, nil
}

func (r *Repository) GetHouseFlats(ctx context.Context, houseID int, isModerator bool) ([]entity.Flat, error) {
	query := r.qb.Select("*").From(flatTable).Where(sq.Eq{"house_id": houseID})

	if !isModerator {
		query = query.Where(sq.Eq{"status": "approved"})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	var rows []flatRow
	err = r.conn.SelectContext(ctx, &rows, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting flats: %w", err)
	}

	result := make([]entity.Flat, 0, len(rows))
	for _, row := range rows {
		result = append(result, entity.Flat{
			ID:         row.ID,
			HouseID:    row.HouseID,
			FlatNumber: row.FlatNumber,
			Price:      row.Price,
			Rooms:      row.Rooms,
			Status:     row.Status,
			CreatedAt:  row.CreatedAt,
			UpdatedAt:  row.UpdatedAt,
		})
	}

	return result, nil
}

func (r *Repository) CreateSubscription(ctx context.Context, houseID int, email string) error {
	insertMap := map[string]any{
		"house_id":   houseID,
		"email":      email,
		"created_at": time.Now(),
	}

	sql, args, err := r.qb.
		Insert(subsTable).
		SetMap(insertMap).
		ToSql()

	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	_, err = r.conn.ExecContext(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, entity.ErrorHouseNotFound) {
			return entity.ErrorHouseNotFound
		}
		return fmt.Errorf("error creating subscription: %w", err)
	}

	return nil
}
