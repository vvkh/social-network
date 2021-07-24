package friendship

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
)

//go:generate mockgen -destination=mocks/usecase.go -package=mocks -source=usecase.go

type UseCase interface {
	ListFriends(ctx context.Context, profileID uint64) ([]entity.Profile, error)
	ListPendingRequests(ctx context.Context, profileID uint64) ([]entity.Profile, error)
	HasPendingRequest(ctx context.Context, profileFromID uint64, profileToID uint64) (bool, error)

	CreateRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error
	AcceptRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error
	DeclineRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error

	StopFriendship(ctx context.Context, profileID uint64, otherProfileID uint64) error
}
