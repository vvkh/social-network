package usecase

import (
	"context"
)

func (u *usecase) CreateUser(ctx context.Context, username string, password string, firstName string, lastName string, age uint8, location string, sex string, about string) (uint64, uint64, error) {
	userID, err := u.repo.CreateUser(ctx, username, password)
	if err != nil {
		return 0, 0, err
	}
	profile, err := u.profilesUC.CreateProfile(ctx, userID, firstName, lastName, age, location, sex, about)
	if err != nil {
		return 0, 0, err
	}
	return userID, profile.ID, nil
}
