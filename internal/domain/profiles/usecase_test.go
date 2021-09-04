package profiles_test

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/vvkh/social-network/internal/config"
	"github.com/vvkh/social-network/internal/db"
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

	conf := config.NewFromEnv()
	appDB, err := db.New(conf.DBUrl)
	require.NoError(t, err)

	profileRepo := profilesRepository.New(appDB)
	require.NoError(t, err)

	profilesUC := profilesUseCase.New(profileRepo)

	repo := usersRepository.New(appDB, conf.BcryptCost)
	require.NoError(t, err)

	uc := usersUseCase.New(profilesUC, repo, "secret")

	ctx := context.Background()
	johnID, johnProfileID, err := uc.CreateUser(ctx, "johndoe_profiles", "123", "john", "doe", 18, "USA", "male", "")
	require.NoError(t, err)
	defer uc.DeleteUser(ctx, johnID) //nolint:errcheck

	topsyID, topsyProfileID, err := uc.CreateUser(ctx, "topsycret_profiles", "123", "topsy", "cret", 19, "USA", "male", "")
	require.NoError(t, err)
	defer uc.DeleteUser(ctx, topsyID) //nolint:errcheck

	profiles, hasMore, err := profilesUC.ListProfiles(ctx, "", "", -1)
	require.NoError(t, err)
	require.False(t, hasMore)

	wantJohnProfile := entity.Profile{
		ID:        johnProfileID,
		UserID:    johnID,
		FirstName: "john",
		LastName:  "doe",
		Age:       18,
		Sex:       "male",
		Location:  "USA",
	}
	wantTopsyProfile := entity.Profile{
		ID:        topsyProfileID,
		UserID:    topsyID,
		FirstName: "topsy",
		LastName:  "cret",
		Age:       19,
		Sex:       "male",
		About:     "",
		Location:  "USA",
	}
	require.Contains(t, profiles, wantJohnProfile)
	require.Contains(t, profiles, wantTopsyProfile)

	johnIndex := 0
	topsyIndex := 0
	for i, profile := range profiles {
		if profile == wantTopsyProfile {
			topsyIndex = i
		}
		if profile == wantJohnProfile {
			johnIndex = i
		}
	}
	require.Equal(t, johnIndex < topsyIndex, johnProfileID < topsyProfileID)
}

func TestSearchProfiles(t *testing.T) {
	if os.Getenv("SKIP_DB_TEST") == "1" {
		t.SkipNow()
	}
	err := godotenv.Load("../../../.env")
	require.NoError(t, err)

	conf := config.NewFromEnv()
	appDB, err := db.New(conf.DBUrl)
	require.NoError(t, err)

	profileRepo := profilesRepository.New(appDB)
	require.NoError(t, err)

	profilesUC := profilesUseCase.New(profileRepo)

	repo := usersRepository.New(appDB, conf.BcryptCost)
	require.NoError(t, err)

	uc := usersUseCase.New(profilesUC, repo, "secret")

	ctx := context.Background()
	johnID, johnProfileID, err := uc.CreateUser(ctx, "johndoe_profiles", "123", "john", "doe", 18, "", "male", "")
	require.NoError(t, err)
	defer uc.DeleteUser(ctx, johnID) //nolint:errcheck

	userID, _, err := uc.CreateUser(ctx, "topsycret_profiles", "123", "topsy", "cret", 18, "", "male", "")
	require.NoError(t, err)
	defer uc.DeleteUser(ctx, userID) //nolint:errcheck

	johnProfile, err := profilesUC.GetByID(ctx, johnProfileID)
	require.NoError(t, err)
	require.Equal(t, 1, len(johnProfile))

	profiles, _, err := profilesUC.ListProfiles(ctx, "john", "doe", -1)
	require.NoError(t, err)
	require.Contains(t, profiles, johnProfile[0])

	profiles, _, err = profilesUC.ListProfiles(ctx, "john", "", -1)
	require.NoError(t, err)
	require.Contains(t, profiles, johnProfile[0])

	profiles, _, err = profilesUC.ListProfiles(ctx, "", "doe", -1)
	require.NoError(t, err)
	require.Contains(t, profiles, johnProfile[0])

	profiles, _, err = profilesUC.ListProfiles(ctx, "jo", "do", -1)
	require.NoError(t, err)
	require.Contains(t, profiles, johnProfile[0])

	profiles, _, err = profilesUC.ListProfiles(ctx, "jo", "", -1)
	require.NoError(t, err)
	require.Contains(t, profiles, johnProfile[0])

	profiles, _, err = profilesUC.ListProfiles(ctx, "", "do", -1)
	require.NoError(t, err)
	require.Contains(t, profiles, johnProfile[0])
}
func TestSearchLimit(t *testing.T) {
	if os.Getenv("SKIP_DB_TEST") == "1" {
		t.SkipNow()
	}
	err := godotenv.Load("../../../.env")
	require.NoError(t, err)

	conf := config.NewFromEnv()
	appDB, err := db.New(conf.DBUrl)
	require.NoError(t, err)

	profileRepo := profilesRepository.New(appDB)
	require.NoError(t, err)

	profilesUC := profilesUseCase.New(profileRepo)

	repo := usersRepository.New(appDB, conf.BcryptCost)
	require.NoError(t, err)

	uc := usersUseCase.New(profilesUC, repo, "secret")

	ctx := context.Background()
	johnID, _, err := uc.CreateUser(ctx, "johndoe_profiles", "123", "john", "doe", 18, "", "male", "")
	require.NoError(t, err)
	defer uc.DeleteUser(ctx, johnID) //nolint:errcheck

	userID, _, err := uc.CreateUser(ctx, "topsycret_profiles", "123", "topsy", "cret", 18, "", "male", "")
	require.NoError(t, err)
	defer uc.DeleteUser(ctx, userID) //nolint:errcheck

	profiles, hasMore, err := profilesUC.ListProfiles(ctx, "", "", -1)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(profiles), 2)
	require.False(t, hasMore)

	profiles, hasMore, err = profilesUC.ListProfiles(ctx, "", "", 1)
	require.NoError(t, err)
	require.Equal(t, 1, len(profiles))
	require.True(t, hasMore)

	profiles, _, err = profilesUC.ListProfiles(ctx, "", "", 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(profiles))

	profiles, _, err = profilesUC.ListProfiles(ctx, "", "", 3)
	require.NoError(t, err)
	require.Equal(t, 2, len(profiles))

	profiles, hasMore, err = profilesUC.ListProfiles(ctx, "", "", 10)
	require.NoError(t, err)
	require.LessOrEqual(t, len(profiles), 10)
	require.False(t, hasMore)
}
