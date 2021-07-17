package repository

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func (r *repo) CreateUser(ctx context.Context, username string, password string) (uint64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	query := `
	INSERT INTO users (username, password) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, query, username, hash)
	if err != nil {
		return 0, err
	}
	userID, err := result.LastInsertId()
	return uint64(userID), err
}
