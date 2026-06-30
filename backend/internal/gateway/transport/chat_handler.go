package transport

import (
	"errors"
	"net/http"
	"strconv"

	"innoconnect/pkg/logger"
	chatpb "innoconnect/pkg/pb/chat"

	"github.com/gin-gonic/gin"
)

type chatMessageRequest struct {
	Text string `json:"text"`
}

func (h *Handler) GetOrCreateRequestChat(c *gin.Context) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	userID, err := extractUserID(c)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid meeting id"})
		return
	}

	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.chatClient.GetMeetingChat(
		c.Request.Context(),
		&chatpb.GetMeetingChatRequest{
			MeetingId: meetingID,
			UserId:    userID,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		return
	}

	userID, err := extractUserID(c)
	if err != nil {
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
	userID, err := extractUserID(c)
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

	userID, err := extractUserID(c)
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

	userID, err := extractUserID(c)
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

func extractUserID(c *gin.Context) (int64, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		logger.Error("missing authorization header")
		return 0, errors.New("missing authorization header")
	}

	// TODO: merge with auth
	return 1, nil
}
