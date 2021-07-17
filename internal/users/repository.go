package users

import "context"

type Repository interface {
	CreateUser(ctx context.Context, username string, password string) (uint64, error)
	Login(ctx context.Context, username string, password string) (uint64, error)
	DeleteUser(ctx context.Context, id uint64) error
}
