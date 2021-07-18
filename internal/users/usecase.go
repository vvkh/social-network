package users

import (
	"context"

	"github.com/vvkh/social-network/internal/users/entity"
)

type UseCase interface {
	CreateUser(ctx context.Context, username string, password string, firstName string, lastName string, age uint8, location string, sex string, about string) (uint64, uint64, error)
	Login(ctx context.Context, username string, password string) (string, error)
	DeleteUser(ctx context.Context, userID uint64) error
	DecodeToken(ctx context.Context, token string) (entity.AccessToken, error)
}
