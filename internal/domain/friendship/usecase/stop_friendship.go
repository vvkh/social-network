package usecase

import (
	"context"
)

func (u *usecase) StopFriendship(ctx context.Context, profileID uint64, otherProfileID uint64) error {
	return u.repo.StopFriendship(ctx, profileID, otherProfileID)
}
