package dto

type Friendship struct {
	RequestedFrom uint64 `db:"requested_from"`
	RequestedTo   uint64 `db:"requested_to"`
	State         string `db:"state"`
}
