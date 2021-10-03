package entity

import "time"

type Chat struct {
	Title               string
	UnreadMessagesCount int64
	ID                  uint64
}

type Message struct {
	Content         string
	AuthorProfileID uint64
	SentAt          time.Time
}
