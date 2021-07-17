package profiles

import (
	"context"

	"github.com/vvkh/social-network/internal/profiles/entity"
)

//go:generate mockgen -destination=mocks/usecase.go -package=mocks -source=usecase.go

type UseCase interface {
	CreateProfile(ctx context.Context, userID uint64, firstName string, lastName string, age uint8, location string, sex string, about string) (entity.Profile, error)
	GetByID(ctx context.Context, id ...uint64) ([]entity.Profile, error)
}
