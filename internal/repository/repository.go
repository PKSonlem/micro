// В слое репозитория мы описываем работу с нашей базой данных
package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	conn *sqlx.DB
	qb   sq.StatementBuilderType
}

func NewRepository(conn *sqlx.DB) *Repository {
	return &Repository{
		conn: conn,
		// инициализируем querybuilder тут, чтобы каждый раз не писать PlaceholderFormat(sq.Dollar)
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
