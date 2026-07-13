package entity

type User struct {
	ID                     int64
	Email                  string
	PasswordHash           string
	Name                   string
	CreatedRequestsCount   int64
	CompletedRequestsCount int64
}
