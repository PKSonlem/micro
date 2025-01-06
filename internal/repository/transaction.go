package repository

import (
	"context"
	"fmt"

	"errors"

	"github.com/jmoiron/sqlx"
)

type txFunc func(tx *sqlx.Tx) error

// Обертка для транзакций, txFunc - это функция внутри которой должны быть вызовы атомарных методов репозитория
func sqlxTransaction(ctx context.Context, db *sqlx.DB, f txFunc) error {
	var txErr error

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	err = f(tx)
	if err != nil {
		txErr = tx.Rollback()
		if txErr != nil {
			return errors.Join(err, txErr)
		}

		return fmt.Errorf("error during transcation: %w", err)
	}

	txErr = tx.Commit()
	if txErr != nil {
		return fmt.Errorf("error commiting transaction: %w", err)
	}

	return nil
}
