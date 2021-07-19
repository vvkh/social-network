package users_test

import (
	"context"
	"os"
	"testing"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"

	"github.com/vvkh/social-network/internal/domain/users"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	profilesRepository "github.com/vvkh/social-network/internal/domain/profiles/repository"
	profilesUseCase "github.com/vvkh/social-network/internal/domain/profiles/usecase"
	"github.com/vvkh/social-network/internal/domain/users/repository"
	"github.com/vvkh/social-network/internal/domain/users/usecase"
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

	uc := usecase.New(profilesUC, repo, "secret")

	ctx := context.Background()

	_, err = uc.Login(ctx, "johndoe1", "topsecret")
	require.Equal(t, users.AuthenticationFailed, err)

	createdUserID, profileID, err := uc.CreateUser(ctx, "johndoe1", "topsecret", "John", "Doe", 18, "USA", "male", "")
	require.NoError(t, err)
	defer uc.DeleteUser(ctx, createdUserID) //nolint:errcheck

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

	rawToken, err := uc.Login(ctx, "johndoe1", "topsecret")
	require.NoError(t, err)

	token, err := uc.DecodeToken(ctx, rawToken)
	require.NoError(t, err)

	loggedInUserID, gotProfileID := token.UserID, token.ProfileID
	require.Equal(t, createdUserID, loggedInUserID)
	require.Equal(t, profileID, gotProfileID)

	_, err = uc.Login(ctx, "johndoe1", "wrongpassword")
	require.Equal(t, users.AuthenticationFailed, err)
}
