package entity

import "time"

type Message struct {
	Content         string
	AuthorProfileID int64
	SentAt          time.Time
}
