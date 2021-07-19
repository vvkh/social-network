package users

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/users/entity"
)

//go:generate mockgen -destination=mocks/usecase.go -package=mocks -source=usecase.go

type UseCase interface {
	CreateUser(ctx context.Context, username string, password string, firstName string, lastName string, age uint8, location string, sex string, about string) (uint64, uint64, error)
	Login(ctx context.Context, username string, password string) (string, error)
	DeleteUser(ctx context.Context, userID uint64) error
	DecodeToken(ctx context.Context, token string) (entity.AccessToken, error)
}
