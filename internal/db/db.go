package db

import "github.com/jmoiron/sqlx"

func New(url string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", url)
	return db, err
}
