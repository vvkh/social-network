package usecase

import "github.com/vvkh/social-network/internal/profiles"

type usecase struct {
	repository profiles.Repository
}

func New(repository profiles.Repository) *usecase {
	return &usecase{
		repository: repository,
	}
}
