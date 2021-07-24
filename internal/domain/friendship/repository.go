package friendship

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/friendship/entity"
)

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go

type Repository interface {
	ListFriends(ctx context.Context, profileID uint64) ([]uint64, error)
	ListPendingRequests(ctx context.Context, profileID uint64) ([]uint64, error)
	GetFriendshipStatus(ctx context.Context, one uint64, other uint64) (entity.FriendshipStatus, error)

	CreateRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error
	AcceptRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error
	DeclineRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error

	StopFriendship(ctx context.Context, profileID uint64, otherProfileID uint64) error
}
