package transport

import (
	"time"

	"innoconnect/internal/chat/entity"
	pb "innoconnect/pkg/pb/chat"
)

func toUser(user entity.User) *pb.User {
	return &pb.User{
		Id:   user.ID,
		Name: user.Name,
	}
}

func toMessage(message entity.Message) *pb.Message {
	return &pb.Message{
		Id:     message.ID,
		Sender: toUser(message.Sender),
		Text:   message.Text,
		SentAt: message.SentAt.Format(time.RFC3339),
	}
}

func toChatResponse(chat entity.Chat) *pb.ChatResponse {
	return &pb.ChatResponse{
		ChatId: chat.ID,
		Type:   toProtoChatType(chat.Type),
		Name: chat.Name,
	}
}

func toGetChatResponse(chat entity.Chat) *pb.GetChatResponse {
	users := make([]*pb.User, 0, len(chat.Participants))

	for _, user := range chat.Participants {
		users = append(users, toUser(user))
	}

	return &pb.GetChatResponse{
		ChatId:      chat.ID,
		ChatName: chat.Name,
		Type:        toProtoChatType(chat.Type),
		Participants: users,
	}
}

func toChatPreview(chat entity.ChatPreview) *pb.ChatPreview {
	users := make([]*pb.User, 0, len(chat.Participants))

	for _, user := range chat.Participants {
		users = append(users, toUser(user))
	}

	var lastMessage *pb.Message
	if chat.LastMessage != nil {
		lastMessage = toMessage(*chat.LastMessage)
	}

	return &pb.ChatPreview{
		ChatId:      chat.ID,
		Name: chat.Name,
		Type:        toProtoChatType(chat.Type),
		Participants: users,
		LastMessage: lastMessage,
	}
}

func toProtoChatType(t entity.ChatType) pb.ChatType {
	switch t {
	case entity.RequestChat:
		return pb.ChatType_REQUEST
	case entity.MeetingChat:
		return pb.ChatType_MEETING
	default:
		return pb.ChatType_CHAT_TYPE_UNSPECIFIED
	}
}