package navbar

type Context struct {
	UnreadMessagesCount            *int64
	PendingFriendshipRequestsCount *int
	SelfProfileID                  *uint64
}
