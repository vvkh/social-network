package entity

const (
	StatePending   = "pending"
	StateAccepted  = "accepted"
	StatesDeclined = "declined"
	StateNone      = "none"
)

type FriendshipStatus struct {
	RequestedFromProfileID uint64
	RequestedToProfileID   uint64
	State                  string
}

func (f FriendshipStatus) IsNone() bool {
	return f.State == StateNone
}

func (f FriendshipStatus) IsPending() bool {
	return f.State == StatePending
}

func (f FriendshipStatus) IsWaitingApprovalFrom(id uint64) bool {
	return f.IsPending() && f.RequestedToProfileID == id
}

func (f FriendshipStatus) IsAccepted() bool {
	return f.State == StateAccepted
}

func (f FriendshipStatus) IsDeclined() bool {
	return f.State == StatesDeclined
}

func (f FriendshipStatus) IsDeclinedBy(id uint64) bool {
	return f.IsDeclined() && f.RequestedToProfileID == id
}
