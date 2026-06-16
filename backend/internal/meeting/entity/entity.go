package entity

import "time"

type Meeting struct {
	ID          int64
	CreatorID   int64
	CreatorName string
	Title       string
	Description string
	Type        string
	Address     *string
	Latitude    *float64
	Longitude   *float64
	MeetingTime time.Time
	MaxPeople   *int32
}
