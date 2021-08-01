package profiles

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
)

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go

type Repository interface {
	CreateProfile(ctx context.Context, profile entity.Profile) (uint64, error)
	GetByID(ctx context.Context, id ...uint64) ([]entity.Profile, error)
	GetByUserID(ctx context.Context, id uint64) ([]entity.Profile, error)
	ListProfiles(ctx context.Context) ([]entity.Profile, error)
	GetByName(ctx context.Context, firstNamePrefix string, lastNamePrefix string) ([]entity.Profile, error)
}
