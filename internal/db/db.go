package db

import (
	_ "github.com/doug-martin/goqu/v9/dialect/mysql" //nolint
	"github.com/jmoiron/sqlx"
)

func New(url string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", url)
	return db, err
}
