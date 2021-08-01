package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
)

func (u *usecase) GetByName(ctx context.Context, firstNamePrefix string, lastNamePrefix string) ([]entity.Profile, error) {
	return u.repository.GetByName(ctx, firstNamePrefix, lastNamePrefix)
}
