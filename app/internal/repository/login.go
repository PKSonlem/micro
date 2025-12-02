package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

type rowUserLogin struct {
	UserId    string    `db:"user_id"`
	UserType  string    `db:"user_type"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

// не знаю зачем написал, ну потом мож пригодиться)
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var res *entity.User
	var err, txErr error

	txErr = sqlxTransaction(ctx, r.conn, func(tx *sqlx.Tx) error {
		res, err = r.getUserByEmail(ctx, email, tx)
		return err
	})
	if txErr != nil {
		return nil, txErr
	}

	return res, nil
}

func (r *Repository) getUserByEmail(ctx context.Context, email string, tx *sqlx.Tx) (*entity.User, error) {
	query := r.qb.Select("*").From(userTable).Where(sq.Eq{"email": email})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error build query: %w", err)
	}

	var row rowUserLogin
	err = tx.GetContext(ctx, &row, sql, args...)
	if err != nil {
		return nil, errors.Join(entity.ErrorLoginUser, err)
	}

	result := &entity.User{
		UserId:       row.UserId,
		UserType:     row.UserType,
		Email:        row.Email,
		PasswordHash: row.Password,
		CreatedAt:    row.CreatedAt,
	}

	return result, nil
}

func (r *Repository) GetUserById(ctx context.Context, userId string) (*entity.User, error) {
	var res *entity.User
	var err, txErr error

	txErr = sqlxTransaction(ctx, r.conn, func(tx *sqlx.Tx) error {
		res, err = r.getUserById(ctx, userId, tx)
		return err
	})
	if txErr != nil {
		return nil, txErr
	}

	return res, nil
}

func (r *Repository) getUserById(ctx context.Context, userId string, tx *sqlx.Tx) (*entity.User, error) {
	query := r.qb.Select("*").From(userTable).Where(sq.Eq{"user_id": userId})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error build query: %w", err)
	}

	var row rowUserLogin
	err = tx.GetContext(ctx, &row, sql, args...)
	if err != nil {
		return nil, errors.Join(entity.ErrorLoginUser, err)
	}

	result := &entity.User{
		UserId:       row.UserId,
		UserType:     row.UserType,
		Email:        row.Email,
		PasswordHash: row.Password,
		CreatedAt:    row.CreatedAt,
	}

	return result, nil
}
