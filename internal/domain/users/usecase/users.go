package usecase

import (
	"github.com/vvkh/social-network/internal/domain/profiles"
	"github.com/vvkh/social-network/internal/domain/users"
)

type usecase struct {
	profilesUC profiles.UseCase
	repo       users.Repository
	jwtSecret  string
}

func New(profilesUC profiles.UseCase, repo users.Repository, jwtSecret string) *usecase {
	return &usecase{
		profilesUC: profilesUC,
		repo:       repo,
		jwtSecret:  jwtSecret,
	}
}
