package chats_test

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/vvkh/social-network/internal/config"
	"github.com/vvkh/social-network/internal/db"
	chatsUseCase "github.com/vvkh/social-network/internal/domain/chats/usecase"
	profilesRepository "github.com/vvkh/social-network/internal/domain/profiles/repository"
	profilesUseCase "github.com/vvkh/social-network/internal/domain/profiles/usecase"
	usersRepository "github.com/vvkh/social-network/internal/domain/users/repository"
	usersUseCase "github.com/vvkh/social-network/internal/domain/users/usecase"
)

func Test_UseCase(t *testing.T) {
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

	usersRepo := usersRepository.New(appDB, conf.BcryptCost)
	require.NoError(t, err)

	usersUC := usersUseCase.New(profilesUC, usersRepo, "secret")

	chatsUC := chatsUseCase.New()

	ctx := context.Background()
	johnID, johnProfileID, err := usersUC.CreateUser(ctx, "johndoe_profiles", "123", "john", "doe", 18, "", "male", "")
	require.NoError(t, err)
	defer usersUC.DeleteUser(ctx, johnID) //nolint:errcheck

	userID, _, err := usersUC.CreateUser(ctx, "topsycret_profiles", "123", "topsy", "cret", 18, "", "male", "")
	require.NoError(t, err)
	defer usersUC.DeleteUser(ctx, userID) //nolint:errcheck

	count, err := chatsUC.GetUnreadMessagesCount(ctx, johnProfileID)
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}
