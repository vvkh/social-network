package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/profiles/entity"
)

func (u *usecase) ListPendingRequests(ctx context.Context, userID uint64) ([]entity.Profile, error) {
	userIDs, err := u.repo.ListPendingRequests(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(userIDs) == 0 {
		return nil, nil
	}

	// TODO: transactions between use cases
	return u.profilesUC.GetByID(ctx, userIDs...)
}
