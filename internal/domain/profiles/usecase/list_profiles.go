package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
)

func (u *usecase) ListProfiles(ctx context.Context) ([]entity.Profile, error) {
	return u.repository.ListProfiles(ctx)
}
