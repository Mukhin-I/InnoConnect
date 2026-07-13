package transport

import (
	"context"
	"net/http"
	"strconv"

	"innoconnect/internal/gateway/entity"
	"innoconnect/internal/gateway/usecase"
	"innoconnect/pkg/logger"
	chatpb "innoconnect/pkg/pb/chat"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type chatMessageRequest struct {
	Text string `json:"text"`
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func (h *Handler) GetOrCreateRequestChat(c *gin.Context) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	userID, _, err := usecase.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.chatClient.GetOrCreateRequestChat(
		c.Request.Context(),
		&chatpb.GetOrCreateRequestChatRequest{
			RequestId: requestID,
			UserId:    userID,
		},
	)
	if err != nil {
		logger.Error("failed to get or create request chat: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get or create request chat"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetMeetingChat(c *gin.Context) {
	meetingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid meeting id"})
		return
	}

	userID, name, err := usecase.GetUserFromToken(c)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.chatClient.GetMeetingChat(
		c.Request.Context(),
		&chatpb.GetMeetingChatRequest{
			MeetingId: meetingID,
			CreatorId:    userID,
			CreatorName: name,
		},
	)
	if err != nil {
		logger.Error("failed to get meeting chat: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get meeting chat"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetChat(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		return
	}

	userID, _, err := usecase.GetUserFromToken(c)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.chatClient.GetChat(
		c.Request.Context(),
		&chatpb.GetChatRequest{
			ChatId: chatID,
			UserId: userID,
		},
	)
	if err != nil {
		logger.Error("failed to get chat: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get chat"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetChats(c *gin.Context) {
	userID, _, err := usecase.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.chatClient.GetChats(
		c.Request.Context(),
		&chatpb.GetChatsRequest{UserId: userID},
	)
	if err != nil {
		logger.Error("failed to get chats: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get chats"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetMessages(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		return
	}

	userID, _, err := usecase.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.chatClient.GetMessages(
		c.Request.Context(),
		&chatpb.GetMessagesRequest{
			ChatId: chatID,
			UserId: userID,
		},
	)
	if err != nil {
		logger.Error("failed to get messages: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) SendMessage(c *gin.Context) {
	logger.Warn("Expired API!!!! Use websocket instead!!!!")
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		return
	}

	var req chatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userID, _, err := usecase.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.chatClient.SendMessage(
		c.Request.Context(),
		&chatpb.SendMessageRequest{
			ChatId: chatID,
			UserId: userID,
			Text:   req.Text,
		},
	)
	if err != nil {
		logger.Error("failed to send message: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) WebSocket(c *gin.Context) {
	userID, _, err := usecase.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	h.wsHub.AddClient(userID, conn)
	defer h.wsHub.RemoveClient(userID)

	for {
		var req entity.WSRequest

		if err := conn.ReadJSON(&req); err != nil {
			logger.Error(err.Error())
			break
		}

		switch req.Type {

		case "send_message":
			h.handleSendMessage(c.Request.Context(), userID, req)

		default:
			logger.Error("unknown websocket message type")
		}
	}
}

// TODO: can be optimized for a one grpc call instead of 2
func (h *Handler) handleSendMessage(
	ctx context.Context,
	userID int64,
	req entity.WSRequest,
) {
	resp, err := h.chatClient.SendMessage(
		ctx,
		&chatpb.SendMessageRequest{
			ChatId: req.ChatID,
			UserId: userID,
			Text:   req.Text,
		},
	)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	participants, err := h.chatClient.GetParticipants(
		ctx,
		&chatpb.GetParticipantsRequest{
			ChatId: req.ChatID,
		},
	)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	wsMessage := entity.WSResponse{
		Type:    "new_message",
		ChatID:  req.ChatID,
		Message: resp,
	}

	for _, participant := range participants.Participants {
		h.wsHub.Send(participant.Id, wsMessage)
	}
}