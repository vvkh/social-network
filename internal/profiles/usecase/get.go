package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/profiles/entity"
)

func (u *usecase) GetByID(ctx context.Context, id ...uint64) ([]entity.Profile, error) {
	return u.repository.GetByID(ctx, id...)
}
