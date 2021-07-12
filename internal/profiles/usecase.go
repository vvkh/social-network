package profiles

import "github.com/vvkh/social-network/internal/profiles/entity"

//go:generate mockgen -destination=mocks/usecase.go -package=mocks -source=usecase.go

type UseCase interface {
	CreateProfile(firstName string, lastName string, age uint8, location string, sex string, about string) (entity.Profile, error)
	GetByID(id uint64) (entity.Profile, error)
}
