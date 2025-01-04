package cmd

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
