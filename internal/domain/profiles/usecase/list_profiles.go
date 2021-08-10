package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
)

func (u *usecase) ListProfiles(ctx context.Context, firstNamePrefix string, lastNamePrefix string, limit int) ([]entity.Profile, bool, error) {
	return u.repository.ListProfiles(ctx, firstNamePrefix, lastNamePrefix, limit)
}
