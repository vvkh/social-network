package repository

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/vvkh/social-network/internal/users"
)

func (r *repo) Login(ctx context.Context, username string, password string) (uint64, error) {
	query := `
    SELECT id, password FROM users WHERE username = ?`
	var user struct {
		ID       uint64 `db:"id""`
		Password []byte `db:"password"`
	}
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		return 0, users.AuthenticationFailed
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return 0, users.AuthenticationFailed
	}
	return user.ID, nil
}
