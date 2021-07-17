package users_test

import (
	"context"
	"os"
	"testing"

	"github.com/vvkh/social-network/internal/profiles/entity"

	"github.com/vvkh/social-network/internal/users"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	profilesRepository "github.com/vvkh/social-network/internal/profiles/repository"
	profilesUseCase "github.com/vvkh/social-network/internal/profiles/usecase"
	"github.com/vvkh/social-network/internal/users/repository"
	"github.com/vvkh/social-network/internal/users/usecase"
)

func TestCreateUserAndLogin(t *testing.T) {
	if os.Getenv("SKIP_DB_TEST") == "1" {
		t.SkipNow()
	}
	err := godotenv.Load("../../.env")
	require.NoError(t, err)

	profileRepo, err := profilesRepository.NewDefault()
	require.NoError(t, err)

	profilesUC := profilesUseCase.New(profileRepo)

	repo, err := repository.NewDefault()
	require.NoError(t, err)

	uc := usecase.New(profilesUC, repo)

	ctx := context.Background()

	_, _, err = uc.Login(ctx, "johndoe", "topsecret")
	require.Equal(t, users.AuthenticationFailed, err)

	createdUserID, profileID, err := uc.CreateUser(ctx, "johndoe", "topsecret", "John", "Doe", 18, "USA", "male", "")
	require.NoError(t, err)
	defer uc.DeleteUser(ctx, createdUserID)

	profile, err := profilesUC.GetByID(ctx, profileID)
	require.NoError(t, err)

	wantProfile := entity.Profile{
		ID:        profileID,
		UserID:    createdUserID,
		FirstName: "John",
		LastName:  "Doe",
		Age:       18,
		Sex:       "male",
		Location:  "USA",
	}
	require.Equal(t, []entity.Profile{wantProfile}, profile)

	loggedInUserID, profileIDs, err := uc.Login(ctx, "johndoe", "topsecret")
	require.NoError(t, err)
	require.Equal(t, createdUserID, loggedInUserID)
	require.Equal(t, []uint64{profileID}, profileIDs)

	_, _, err = uc.Login(ctx, "johndoe", "wrongpassword")
	require.Equal(t, users.AuthenticationFailed, err)
}
