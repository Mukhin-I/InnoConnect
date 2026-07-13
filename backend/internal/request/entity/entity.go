package entity

import "time"

type Request struct {
	ID               int64
	CreatorID        int64
	CreatorName      string
	Title            string
	Description      string
	RequesterAddress string
	Type             string
	Status string
	Deadline         time.Time
}
