package friendship_test

import (
	"context"
	"os"
	"testing"

	"github.com/vvkh/social-network/internal/profiles/entity"

	"github.com/stretchr/testify/require"

	"github.com/joho/godotenv"

	"github.com/vvkh/social-network/internal/friendship/repository"
	"github.com/vvkh/social-network/internal/friendship/usecase"
	profilesRepository "github.com/vvkh/social-network/internal/profiles/repository"
	profilesUseCase "github.com/vvkh/social-network/internal/profiles/usecase"
)

func TestAcceptFriendshipRequest(t *testing.T) {
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

	uc := usecase.New(repo, profilesUC)

	ctx := context.Background()

	john, err := profilesUC.CreateProfile(ctx, "John", "Doe", 18, "", "male", "")
	require.NoError(t, err)

	topsy, err := profilesUC.CreateProfile(ctx, "Topsy", "Cret", 18, "", "male", "")
	require.NoError(t, err)

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
	err := godotenv.Load("../../.env")
	require.NoError(t, err)

	profileRepo, err := profilesRepository.NewDefault()
	require.NoError(t, err)

	profilesUC := profilesUseCase.New(profileRepo)

	repo, err := repository.NewDefault()
	require.NoError(t, err)

	uc := usecase.New(repo, profilesUC)

	ctx := context.Background()

	john, err := profilesUC.CreateProfile(ctx, "John", "Doe", 18, "", "male", "")
	require.NoError(t, err)

	topsy, err := profilesUC.CreateProfile(ctx, "Topsy", "Cret", 18, "", "male", "")
	require.NoError(t, err)

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
}

func TestStopFriendship(t *testing.T) {
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

	uc := usecase.New(repo, profilesUC)

	ctx := context.Background()

	john, err := profilesUC.CreateProfile(ctx, "John", "Doe", 18, "", "male", "")
	require.NoError(t, err)

	topsy, err := profilesUC.CreateProfile(ctx, "Topsy", "Cret", 18, "", "male", "")
	require.NoError(t, err)

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
