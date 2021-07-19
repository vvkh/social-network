package friendship

import (
	"context"
)

//go:generate mockgen -destination=mocks/repository.go -package=mocks -source=repository.go

type Repository interface {
	ListFriends(ctx context.Context, profileID uint64) ([]uint64, error)
	ListPendingRequests(ctx context.Context, profileID uint64) ([]uint64, error)

	CreateRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error
	AcceptRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error
	DeclineRequest(ctx context.Context, profileFromID uint64, profileToID uint64) error

	StopFriendship(ctx context.Context, profileID uint64, otherProfileID uint64) error
}
