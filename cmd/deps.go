package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewSqlxConn(configuration *postgresConfiguration) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", configuration.GetConnectionString())
	if err != nil {
		return nil, errors.Wrap(err, "cant connect to db")
	}

	db.SetMaxIdleConns(configuration.GetMaxIdleConns())
	db.SetMaxOpenConns(configuration.GetMaxOpenConns())

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "cant ping db")
	}

	return db, nil
}

type Logger struct {
	*slog.Logger
}

func NewLogger() *Logger {
	return &Logger{slog.New(slog.NewJSONHandler(os.Stdout, nil))}
}

func (l *Logger) WithTag(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, "tag", tag)
}

func (l *Logger) Info(ctx context.Context, message string, args ...any) {
	tag := ctx.Value("tag").(string)

	l.Logger.InfoContext(ctx, fmt.Sprintf("tag: %s, message: %s", tag, message), args...)
}

func (l *Logger) Error(ctx context.Context, err error, args ...any) {
	l.Logger.ErrorContext(ctx, err.Error(), args...)
}
