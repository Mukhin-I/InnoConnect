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

// GetOrCreateRequestChat godoc
// @Summary Get or create a request chat
// @Description Returns the chat for a request, creating it if needed
// @Tags chats
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param id path int true "Request ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /requests/{id}/chat [post]
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

// GetMeetingChat godoc
// @Summary Get a meeting chat
// @Description Returns the chat associated with a meeting
// @Tags chats
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param id path int true "Meeting ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /meetings/{id}/chat [get]
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
			MeetingId:   meetingID,
			CreatorId:   userID,
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

// GetChat godoc
// @Summary Get a chat by ID
// @Description Returns detailed information about a specific chat
// @Tags chats
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param chat_id path int true "Chat ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chats/{chat_id} [get]
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

// GetChats godoc
// @Summary Get all chats
// @Description Returns all chats available to the authenticated user
// @Tags chats
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chats [get]
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

// GetMessages godoc
// @Summary Get chat messages
// @Description Returns messages for a chat
// @Tags chats
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param chat_id path int true "Chat ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chats/{chat_id}/messages [get]
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

// SendMessage godoc
// @Summary Send a chat message
// @Description Sends a message to a chat using the REST endpoint
// @Tags chats
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param chat_id path int true "Chat ID"
// @Param request body chatMessageRequest true "Message payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chats/{chat_id}/messages [post]
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

// WebSocket godoc
// @Summary Open a chat websocket
// @Description Opens a websocket connection for real-time chat communication
// @Tags chats
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Success 101 {string} string
// @Failure 401 {object} map[string]interface{}
// @Router /ws [get]
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
