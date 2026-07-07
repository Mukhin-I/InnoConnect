package entity

type ChatType int

const (
	RequestChat ChatType = iota + 1
	MeetingChat
)

type Chat struct {
	ID           int64
	Type         ChatType
	RelatedID int64
	Participants []User
}

type ChatPreview struct {
	ID           int64
	Type         ChatType
	Title        string
	Participants []User
	LastMessage  *Message
}

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