package repository

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db         *sqlx.DB
	bcryptCost int
}

func New(db *sqlx.DB, bcryptCost int) *repo {
	return &repo{
		db:         db,
		bcryptCost: bcryptCost,
	}
}
