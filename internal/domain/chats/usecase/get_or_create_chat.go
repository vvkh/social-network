package usecase

import (
	"context"
)

func (u *uc) GetOrCreateChat(ctx context.Context, oneProfileID uint64, otherProfileID uint64) (uint64, error) {
	return u.repo.GetOrCreateChat(ctx, "title", oneProfileID, otherProfileID)
}
