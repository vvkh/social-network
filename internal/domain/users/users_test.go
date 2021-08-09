package users_test

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
	"github.com/vvkh/social-network/internal/domain/users"
	"github.com/vvkh/social-network/internal/domain/users/repository"
	"github.com/vvkh/social-network/internal/domain/users/usecase"
)

func TestCreateUserAndLogin(t *testing.T) {
	if os.Getenv("SKIP_DB_TEST") == "1" {
		t.SkipNow()
	}
	err := godotenv.Load("../../../.env")
	require.NoError(t, err)

	conf := config.NewFromEnv()
	appDB, err := db.New(conf.DBUrl)
	require.NoError(t, err)

	profileRepo := profilesRepository.New(appDB)
	profilesUC := profilesUseCase.New(profileRepo)

	repo := repository.New(appDB, conf.BcryptCost)
	uc := usecase.New(profilesUC, repo, "secret")

	ctx := context.Background()

	_, err = uc.Login(ctx, "johndoe_users", "topsecret")
	require.Equal(t, users.AuthenticationFailed, err)

	createdUserID, profileID, err := uc.CreateUser(ctx, "johndoe_users", "topsecret", "John", "Doe", 18, "USA", "male", "")
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

	rawToken, err := uc.Login(ctx, "johndoe_users", "topsecret")
	require.NoError(t, err)

	token, err := uc.DecodeToken(ctx, rawToken)
	require.NoError(t, err)

	loggedInUserID, gotProfileID := token.UserID, token.ProfileID
	require.Equal(t, createdUserID, loggedInUserID)
	require.Equal(t, profileID, gotProfileID)

	_, err = uc.Login(ctx, "johndoe_users", "wrongpassword")
	require.Equal(t, users.AuthenticationFailed, err)
}
