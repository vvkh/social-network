package usecase

import (
	"os"

	"github.com/vvkh/social-network/internal/domain/profiles"
	"github.com/vvkh/social-network/internal/domain/users"
)

type usecase struct {
	profilesUC profiles.UseCase
	repo       users.Repository
	jwtSecret  string
}

func NewFromEnv(profilesUC profiles.UseCase, repo users.Repository) *usecase {
	jwtSecret := os.Getenv("AUTH_SECRET")
	return New(profilesUC, repo, jwtSecret)
}

func New(profilesUC profiles.UseCase, repo users.Repository, jwtSecret string) *usecase {
	return &usecase{
		profilesUC: profilesUC,
		repo:       repo,
		jwtSecret:  jwtSecret,
	}
}
