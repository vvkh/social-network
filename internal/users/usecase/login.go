package usecase

import (
	"context"
)

func (u *usecase) Login(ctx context.Context, username string, password string) (uint64, []uint64, error) {
	userID, err := u.repo.Login(ctx, username, password)
	if err != nil {
		return 0, nil, err
	}
	profiles, err := u.profilesUC.GetByUserID(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	profileIDs := make([]uint64, 0, len(profiles))
	for _, profile := range profiles {
		profileIDs = append(profileIDs, profile.ID)
	}
	return userID, profileIDs, nil
}
