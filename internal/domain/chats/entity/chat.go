package entity

import "time"

type Chat struct {
	Title               string
	UnreadMessagesCount int64
	ID                  int64
}

type Message struct {
	Content         string
	AuthorProfileID int64
	SentAt          time.Time
}
