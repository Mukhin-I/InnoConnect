package entity

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