package entity

// Request payload for creating a meeting
type CreateMeetingRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Address     *string  `json:"address,omitempty"`
	Type        string   `json:"type"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
	MeetingTime string   `json:"meeting_time"`
	MaxPeople   *int32   `json:"max_people,omitempty"`
}

// Meeting preview without detailed info
type MeetingShort struct {
	ID        int64    `json:"id"`
	Address   *string  `json:"address,omitempty"`
	Type      string   `json:"type"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// Meeting participant
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Complete meeting with creator, participants and all fields
type MeetingFull struct {
	ID int64 `json:"id"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Creator      User   `json:"creator"`
	Participants []User `json:"participants"`

	Address *string `json:"address,omitempty"`

	Type string `json:"type"`

	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`

	MeetingTime string `json:"meeting_time"`

	CurrentPeople int32 `json:"current_people"`

	MaxPeople *int32 `json:"max_people,omitempty"`
}

// Response wrapper for meetings list
type GetMeetingsResponse struct {
	Meetings []MeetingShort `json:"meetings"`
}
