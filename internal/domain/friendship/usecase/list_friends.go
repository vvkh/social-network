package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
)

func (u *usecase) ListFriends(ctx context.Context, profileID uint64) ([]entity.Profile, error) {
	friendIDs, err := u.repo.ListFriends(ctx, profileID)
	if err != nil {
		return nil, err
	}
	if len(friendIDs) == 0 {
		return nil, nil
	}
	return u.profilesUC.GetByID(ctx, friendIDs...)
}
