package entity

import "time"

type Message struct {
	ID     int64
	Sender User
	Text   string
	SentAt time.Time
}