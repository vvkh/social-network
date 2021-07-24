package friendship

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/friendship/entity"
	profileEntity "github.com/vvkh/social-network/internal/domain/profiles/entity"
)

//go:generate mockgen -destination=mocks/usecase.go -package=mocks -source=usecase.go

type UseCase interface {
	ListFriends(ctx context.Context, profileID uint64) ([]profileEntity.Profile, error)
	ListPendingRequests(ctx context.Context, profileID uint64) ([]profileEntity.Profile, error)
	GetFriendshipStatus(ctx context.Context, one uint64, other uint64) (entity.FriendshipStatus, error)

	CreateRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error
	AcceptRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error
	DeclineRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error

	StopFriendship(ctx context.Context, profileID uint64, otherProfileID uint64) error
}
