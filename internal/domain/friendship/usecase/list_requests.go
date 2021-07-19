package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
)

func (u *usecase) ListPendingRequests(ctx context.Context, profileID uint64) ([]entity.Profile, error) {
	profileIDs, err := u.repo.ListPendingRequests(ctx, profileID)
	if err != nil {
		return nil, err
	}
	if len(profileIDs) == 0 {
		return nil, nil
	}

	// TODO: transactions between use cases
	return u.profilesUC.GetByID(ctx, profileIDs...)
}
