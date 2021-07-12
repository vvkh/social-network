package profiles

import "github.com/vvkh/social-network/internal/profiles/entity"

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go

type Repository interface {
	CreateProfile(profile entity.Profile) (uint64, error)
	GetByID(id uint64) (entity.Profile, error)
}
