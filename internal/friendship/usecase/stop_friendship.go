package usecase

import (
	"context"
)

func (u *usecase) StopFriendship(ctx context.Context, userID uint64, otherUserID uint64) error {
	return u.repo.StopFriendship(ctx, userID, otherUserID)
}
