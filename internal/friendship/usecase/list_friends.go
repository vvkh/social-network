package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/profiles/entity"
)

func (u *usecase) ListFriends(ctx context.Context, userID uint64) ([]entity.Profile, error) {
	friendIDs, err := u.repo.ListFriends(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(friendIDs) == 0 {
		return nil, nil
	}
	return u.profilesUC.GetByID(ctx, friendIDs...)
}
