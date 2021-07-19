package usecase

import (
	"github.com/vvkh/social-network/internal/domain/friendship"
	"github.com/vvkh/social-network/internal/domain/profiles"
)

type usecase struct {
	repo       friendship.Repository
	profilesUC profiles.UseCase
}

func New(repo friendship.Repository, profilesUC profiles.UseCase) *usecase {
	return &usecase{
		repo:       repo,
		profilesUC: profilesUC,
	}
}
