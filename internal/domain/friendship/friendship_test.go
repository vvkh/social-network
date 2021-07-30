package friendship_test

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/vvkh/social-network/internal/config"
	"github.com/vvkh/social-network/internal/db"
	"github.com/vvkh/social-network/internal/domain/friendship/repository"
	"github.com/vvkh/social-network/internal/domain/friendship/usecase"
	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	profilesRepository "github.com/vvkh/social-network/internal/domain/profiles/repository"
	profilesUseCase "github.com/vvkh/social-network/internal/domain/profiles/usecase"
	usersRepository "github.com/vvkh/social-network/internal/domain/users/repository"
	usersUseCase "github.com/vvkh/social-network/internal/domain/users/usecase"
)

func TestAcceptFriendshipRequest(t *testing.T) {
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

	usersRepo := usersRepository.New(appDB)
	usersUC := usersUseCase.New(profilesUC, usersRepo, "secret")

	repo := repository.New(appDB)
	uc := usecase.New(repo, profilesUC)

	ctx := context.Background()

	johnUserID, johnProfileID, err := usersUC.CreateUser(ctx, "johndoe_friendship", "topsecret", "John", "Doe", 18, "", "male", "")
	require.NoError(t, err)
	defer usersUC.DeleteUser(ctx, johnUserID) //nolint:errcheck

	profiles, err := profilesUC.GetByID(ctx, johnProfileID)
	require.NoError(t, err)
	john := profiles[0]

	topsyUserID, topsyProfileID, err := usersUC.CreateUser(ctx, "topsycret", "topsecret", "Topsy", "Cret", 18, "", "male", "")
	require.NoError(t, err)
	defer usersUC.DeleteUser(ctx, topsyUserID) //nolint:errcheck

	profiles, err = profilesUC.GetByID(ctx, topsyProfileID)
	require.NoError(t, err)
	topsy := profiles[0]

	friends, err := uc.ListFriends(ctx, topsy.ID)
	require.NoError(t, err)
	require.Empty(t, friends)

	friends, err = uc.ListFriends(ctx, john.ID)
	require.NoError(t, err)
	require.Empty(t, friends)

	err = uc.CreateRequest(ctx, john.ID, topsy.ID)
	require.NoError(t, err)

	requests, err := uc.ListPendingRequests(ctx, topsy.ID)
	require.NoError(t, err)
	require.Equal(t, []entity.Profile{john}, requests)

	requests, err = uc.ListPendingRequests(ctx, john.ID)
	require.NoError(t, err)
	require.Empty(t, requests)

	err = uc.AcceptRequest(ctx, john.ID, topsy.ID)
	require.NoError(t, err)

	topsyFriends, err := uc.ListFriends(ctx, topsy.ID)
	require.NoError(t, err)
	require.Equal(t, []entity.Profile{john}, topsyFriends)

	johnFriends, err := uc.ListFriends(ctx, john.ID)
	require.NoError(t, err)
	require.Equal(t, []entity.Profile{topsy}, johnFriends)
}

func TestDeclineFriendshipRequest(t *testing.T) {
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

	usersRepo := usersRepository.New(appDB)
	usersUC := usersUseCase.New(profilesUC, usersRepo, "secret")

	repo := repository.New(appDB)
	uc := usecase.New(repo, profilesUC)

	ctx := context.Background()

	johnUserID, johnProfileID, err := usersUC.CreateUser(ctx, "johndoe_friendship", "topsecret", "John", "Doe", 18, "", "male", "")
	require.NoError(t, err)
	defer usersUC.DeleteUser(ctx, johnUserID) //nolint:errcheck

	profiles, err := profilesUC.GetByID(ctx, johnProfileID)
	require.NoError(t, err)
	john := profiles[0]

	topsyUserID, topsyProfileID, err := usersUC.CreateUser(ctx, "topsycret", "topsecret", "Topsy", "Cret", 18, "", "male", "")
	require.NoError(t, err)
	defer usersUC.DeleteUser(ctx, topsyUserID) //nolint:errcheck

	profiles, err = profilesUC.GetByID(ctx, topsyProfileID)
	require.NoError(t, err)
	topsy := profiles[0]

	friends, err := uc.ListFriends(ctx, topsy.ID)
	require.NoError(t, err)
	require.Empty(t, friends)

	friends, err = uc.ListFriends(ctx, john.ID)
	require.NoError(t, err)
	require.Empty(t, friends)

	err = uc.CreateRequest(ctx, john.ID, topsy.ID)
	require.NoError(t, err)

	requests, err := uc.ListPendingRequests(ctx, topsy.ID)
	require.NoError(t, err)
	require.Equal(t, []entity.Profile{john}, requests)

	requests, err = uc.ListPendingRequests(ctx, john.ID)
	require.NoError(t, err)
	require.Empty(t, requests)

	err = uc.DeclineRequest(ctx, john.ID, topsy.ID)
	require.NoError(t, err)

	topsyFriends, err := uc.ListFriends(ctx, topsy.ID)
	require.NoError(t, err)
	require.Empty(t, topsyFriends)

	johnFriends, err := uc.ListFriends(ctx, john.ID)
	require.NoError(t, err)
	require.Empty(t, johnFriends)

	requests, err = uc.ListPendingRequests(ctx, topsy.ID)
	require.NoError(t, err)
	require.Empty(t, requests)

	err = uc.AcceptRequest(ctx, john.ID, topsy.ID)
	require.NoError(t, err)

	johnFriends, err = uc.ListFriends(ctx, john.ID)
	require.NoError(t, err)
	require.Contains(t, johnFriends, topsy)
}

func TestStopFriendship(t *testing.T) {
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

	usersRepo := usersRepository.New(appDB)
	usersUC := usersUseCase.New(profilesUC, usersRepo, "secret")

	repo := repository.New(appDB)
	uc := usecase.New(repo, profilesUC)

	ctx := context.Background()

	johnUserID, johnProfileID, err := usersUC.CreateUser(ctx, "johndoe_friendship", "topsecret", "John", "Doe", 18, "", "male", "")
	require.NoError(t, err)
	defer usersUC.DeleteUser(ctx, johnUserID) //nolint:errcheck

	profiles, err := profilesUC.GetByID(ctx, johnProfileID)
	require.NoError(t, err)
	john := profiles[0]

	topsyUserID, topsyProfileID, err := usersUC.CreateUser(ctx, "topsycret", "topsecret", "Topsy", "Cret", 18, "", "male", "")
	require.NoError(t, err)
	defer usersUC.DeleteUser(ctx, topsyUserID) //nolint:errcheck

	profiles, err = profilesUC.GetByID(ctx, topsyProfileID)
	require.NoError(t, err)
	topsy := profiles[0]

	err = uc.CreateRequest(ctx, john.ID, topsy.ID)
	require.NoError(t, err)

	err = uc.AcceptRequest(ctx, john.ID, topsy.ID)
	require.NoError(t, err)

	err = uc.StopFriendship(ctx, topsy.ID, john.ID) // note that order might be different
	require.NoError(t, err)

	friends, err := uc.ListFriends(ctx, topsy.ID)
	require.NoError(t, err)
	require.Empty(t, friends)

	friends, err = uc.ListFriends(ctx, john.ID)
	require.NoError(t, err)
	require.Empty(t, friends)
}

func TestGetFriendshipStatus(t *testing.T) {
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

	usersRepo := usersRepository.New(appDB)
	usersUC := usersUseCase.New(profilesUC, usersRepo, conf.AuthSecret)

	repo := repository.New(appDB)
	uc := usecase.New(repo, profilesUC)

	ctx := context.Background()

	johnUserID, johnProfileID, err := usersUC.CreateUser(ctx, "johndoe_friendship", "topsecret", "John", "Doe", 18, "", "male", "")
	require.NoError(t, err)
	defer usersUC.DeleteUser(ctx, johnUserID) //nolint:errcheck

	topsyUserID, topsyProfileID, err := usersUC.CreateUser(ctx, "topsycret", "topsecret", "Topsy", "Cret", 18, "", "male", "")
	require.NoError(t, err)
	defer usersUC.DeleteUser(ctx, topsyUserID) //nolint:errcheck

	status, err := uc.GetFriendshipStatus(ctx, johnProfileID, topsyProfileID)
	require.NoError(t, err)
	require.True(t, status.IsNone())

	err = uc.CreateRequest(ctx, johnProfileID, topsyProfileID)
	require.NoError(t, err)

	status, err = uc.GetFriendshipStatus(ctx, johnProfileID, topsyProfileID)
	require.NoError(t, err)
	require.True(t, status.IsPending())
	require.True(t, status.IsWaitingApprovalFrom(topsyProfileID))

	err = uc.AcceptRequest(ctx, johnProfileID, topsyProfileID)
	require.NoError(t, err)

	status, err = uc.GetFriendshipStatus(ctx, johnProfileID, topsyProfileID)
	require.NoError(t, err)
	require.True(t, status.IsAccepted())
}
