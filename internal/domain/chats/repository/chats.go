package repository

import "github.com/jmoiron/sqlx"

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repo {
	return &repo{
		db: db,
	}
}
