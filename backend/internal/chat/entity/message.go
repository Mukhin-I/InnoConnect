package entity

import "time"

// Message represents a single chat message sent by a user
type Message struct {
	ID     int64
	Sender User
	Text   string
	SentAt time.Time
}
