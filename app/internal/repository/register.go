package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/PKSonlem/micro/internal/entity"
)

const userTable = "users"

type rowUserRegister struct {
	UserId    string    `db:"user_id"`
	UserType  string    `db:"user_type"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *Repository) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	var res *entity.User
	var err, txErr error

	txErr = sqlxTransaction(ctx, r.conn, func(tx *sqlx.Tx) error {
		res, err = r.createUserTx(ctx, user, tx)
		return err
	})
	if txErr != nil {
		return nil, txErr
	}

	return res, nil
}

func (r *Repository) createUserTx(ctx context.Context, user entity.User, tx *sqlx.Tx) (*entity.User, error) {
	insertMap := map[string]any{
		"user_type":  user.UserType,
		"email":      user.Email,
		"password":   user.PasswordHash,
		"created_at": time.Now(),
	}

	sql, args, err := r.qb.Insert(userTable).SetMap(insertMap).Suffix("RETURNING *").ToSql()

	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	var row rowUserRegister
	err = tx.GetContext(ctx, &row, sql, args...)
	if err != nil {
		return nil, errors.Join(entity.ErrorCreatingUser, err)
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
