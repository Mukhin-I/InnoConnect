package entity

type ChatType int

const (
	RequestChat ChatType = iota + 1
	MeetingChat
)

// Chat represents a conversation between users
type Chat struct {
	ID           int64
	Name         string
	Type         ChatType
	RelatedID    int64
	Participants []User
}

// ChatPreview includes last message for display
type ChatPreview struct {
	ID           int64
	Name         string
	Type         ChatType
	Title        string
	Participants []User
	LastMessage  *Message
}

// Convert string to ChatType enum
func ChatTypeFromString(s string) ChatType {
	switch s {
	case "REQUEST":
		return RequestChat
	case "MEETING":
		return MeetingChat
	default:
		return RequestChat // safe fallback for MVP
	}
}
