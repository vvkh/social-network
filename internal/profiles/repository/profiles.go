package repository

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewDefault() (*repo, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error while parsing env: %w", err)
	}

	url := os.Getenv("DB_URL")
	return New(url)
}

func New(url string) (*repo, error) {
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	return &repo{
		db: db,
	}, nil
}
