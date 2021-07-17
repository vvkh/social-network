package usecase

import (
	"context"
	"fmt"

	"github.com/vvkh/social-network/internal/profiles/entity"
)

func (u *usecase) GetByID(ctx context.Context, id uint64) (entity.Profile, error) {
	profile, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return entity.Profile{}, fmt.Errorf("repository.GetByID error: %w", err)
	}
	return profile, nil
}
