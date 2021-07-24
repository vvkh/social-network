package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/friendship/entity"
)

func (u *usecase) GetFriendshipStatus(ctx context.Context, one uint64, other uint64) (entity.FriendshipStatus, error) {
	status, err := u.repo.GetFriendshipStatus(ctx, one, other)
	return status, err
}
