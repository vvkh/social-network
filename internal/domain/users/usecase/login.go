package usecase

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/vvkh/social-network/internal/domain/users"
	"github.com/vvkh/social-network/internal/domain/users/entity"
)

func (u *usecase) Login(ctx context.Context, username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", users.EmptyCredentials
	}
	userID, err := u.repo.Login(ctx, username, password)
	if err != nil {
		return "", err
	}
	profiles, err := u.profilesUC.GetByUserID(ctx, userID)
	if err != nil {
		return "", err
	}
	if len(profiles) == 0 {
		return "", fmt.Errorf("something went wrong, user %d exists but doesn't have any profiles", userID)
	}
	profileIDs := make([]uint64, 0, len(profiles))
	for _, profile := range profiles {
		profileIDs = append(profileIDs, profile.ID)
	}
	sort.Slice(profileIDs, func(i, j int) bool {
		return profileIDs[i] < profileIDs[j]
	})

	token := entity.AccessToken{
		UserID:    userID,
		ProfileID: profileIDs[0], // functionality for multiple profiles is undefined, so just use the first one
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	}

	return token.ToString(u.jwtSecret)
}
