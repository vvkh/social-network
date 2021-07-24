package dto

import "github.com/vvkh/social-network/internal/domain/friendship/entity"

type Friendship struct {
	RequestedFrom uint64 `db:"requested_from"`
	RequestedTo   uint64 `db:"requested_to"`
	State         string `db:"state"`
}

func ToModel(friendship Friendship) entity.FriendshipStatus {
	if friendship.State == "created" {
		return entity.FriendshipStatus{
			State:                  entity.StatePending,
			RequestedFromProfileID: friendship.RequestedFrom,
			RequestedToProfileID:   friendship.RequestedTo,
		}
	}
	if friendship.State == "accepted" {
		return entity.FriendshipStatus{
			State:                  entity.StateAccepted,
			RequestedFromProfileID: friendship.RequestedFrom,
			RequestedToProfileID:   friendship.RequestedTo,
		}
	}
	if friendship.State == "declined" {
		return entity.FriendshipStatus{
			State:                  entity.StatesDeclined,
			RequestedFromProfileID: friendship.RequestedFrom,
			RequestedToProfileID:   friendship.RequestedTo,
		}
	}
	return entity.FriendshipStatus{State: entity.StateNone}
}
