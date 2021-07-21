package profiles_test

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	profilesRepository "github.com/vvkh/social-network/internal/domain/profiles/repository"
	profilesUseCase "github.com/vvkh/social-network/internal/domain/profiles/usecase"
	usersRepository "github.com/vvkh/social-network/internal/domain/users/repository"
	usersUseCase "github.com/vvkh/social-network/internal/domain/users/usecase"
)

func TestProfiles(t *testing.T) {
	if os.Getenv("SKIP_DB_TEST") == "1" {
		t.SkipNow()
	}
	err := godotenv.Load("../../../.env")
	require.NoError(t, err)

	profileRepo, err := profilesRepository.NewDefault()
	require.NoError(t, err)

	profilesUC := profilesUseCase.New(profileRepo)

	repo, err := usersRepository.NewDefault()
	require.NoError(t, err)

	uc := usersUseCase.New(profilesUC, repo, "secret")

	ctx := context.Background()
	johnID, johnProfileID, err := uc.CreateUser(ctx, "johndoe", "123", "john", "doe", 18, "USA", "male", "")
	require.NoError(t, err)
	defer uc.DeleteUser(ctx, johnID) //nolint:errcheck

	topsyID, topsyProfileID, err := uc.CreateUser(ctx, "topsycret", "123", "topsy", "cret", 19, "USA", "male", "")
	require.NoError(t, err)
	defer uc.DeleteUser(ctx, topsyID) //nolint:errcheck

	profiles, err := profilesUC.ListProfiles(ctx)
	wantProfiles := []entity.Profile{
		{
			ID:        johnProfileID,
			UserID:    johnID,
			FirstName: "john",
			LastName:  "doe",
			Age:       18,
			Sex:       "male",
			Location:  "USA",
		},
		{
			ID:        topsyProfileID,
			UserID:    topsyID,
			FirstName: "topsy",
			LastName:  "cret",
			Age:       19,
			Sex:       "male",
			About:     "",
			Location:  "USA",
		},
	}
	require.Equal(t, wantProfiles, profiles)

}
