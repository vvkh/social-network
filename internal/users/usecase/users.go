package usecase

import (
	"github.com/vvkh/social-network/internal/profiles"
	"github.com/vvkh/social-network/internal/users"
)

type usecase struct {
	profilesUC profiles.UseCase
	repo       users.Repository
}

func New(profilesUC profiles.UseCase, repo users.Repository) *usecase {
	return &usecase{
		profilesUC: profilesUC,
		repo:       repo,
	}
}
