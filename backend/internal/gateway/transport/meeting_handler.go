package transport

import (
	"net/http"
	"strconv"

	"innoconnect/internal/gateway/entity"
	"innoconnect/internal/gateway/usecase"
	"innoconnect/pkg/logger"
	"innoconnect/pkg/pb/chat"
	"innoconnect/pkg/pb/meeting"

	"github.com/gin-gonic/gin"
)

// CreateMeeting godoc
// @Summary Create a meeting
// @Description Creates a new meeting
// @Tags meetings
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param request body entity.CreateMeetingRequest true "Meeting data"
// @Success 201 {object} entity.MeetingFull
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /meetings [post]
func (h *Handler) CreateMeeting(c *gin.Context) {
	var req entity.CreateMeetingRequest

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		logger.Error("missing authorization header")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Invalid JSON": err.Error(),
		})
		logger.Error("Invalid JSON in meeting request")
		return
	}

	logger.Info("Sending gRPC request to create meeting")

	userID, name, err := usecase.GetUserFromToken(c)
	logger.Info("Username " + name)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	meeting, err := h.meetingClient.CreateMeeting(
		c.Request.Context(),
		&meeting.CreateMeetingRequest{
			CreatorId: userID,
			CreatorName: name,
			Title: req.Title,
			Description: req.Description,
			Address: req.Address,
			Latitude: req.Latitude,
			Longitude: req.Longitude,
			Type: req.Type,
			MeetingTime: req.MeetingTime,
			MaxPeople: req.MaxPeople,
		},
	)

	_, err = h.chatClient.CreateMeetingChat(
		c.Request.Context(),
		&chat.CreateMeetingChatRequest{
			MeetingId: meeting.Id,
			ChatName: req.Title,
			CreatorId: userID,
			CreatorName: name,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logger.Error("Falied go get answer from meeting service via gRPC" + err.Error())
		return
	}

	c.JSON(http.StatusCreated, meeting)
}

// GetMeetings godoc
// @Summary Get all meetings
// @Description Returns a list of meetings
// @Tags meetings
// @Produce json
// @Success 200 {object} entity.GetMeetingsResponse
// @Failure 500 "Internal Server Error"
// @Router /meetings [get]
func (h *Handler) GetMeetings(c *gin.Context) {
	logger.Info("Getting meetings from meeting service")
	resp, err := h.meetingClient.GetMeetings(
		c.Request.Context(),
		&meeting.GetMeetingsRequest{},
	)

	if err != nil {
		logger.Error(
			"Failed to get meetings from meeting service: " +
				err.Error(),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get meetings" + err.Error(),
		})
		return
	}

	var meetings []entity.MeetingShort

	for _, m := range resp.Meetings {
		meetings = append(meetings, entity.MeetingShort{
			ID: m.Id,
			Address: m.Address,
			Type: m.Type,
			Latitude: m.Latitude,
			Longitude: m.Longitude,
		})
	}

	c.JSON(http.StatusOK, entity.GetMeetingsResponse{
		Meetings: meetings,
	})
}

// GetMeeting godoc
// @Summary Get meeting by ID
// @Description Returns detailed information about a meeting
// @Tags meetings
// @Produce json
// @Param id path int true "Meeting ID"
// @Success 200 {object} entity.MeetingFull
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /meetings/{id} [get]
func (h *Handler) GetMeeting(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	logger.Info("Getting meeting by id: " + strconv.FormatInt(id, 10))

	resp, err := h.meetingClient.GetMeeting(
		c.Request.Context(),
		&meeting.GetMeetingRequest{
			Id: id,
		},
	)

	if err != nil {
		logger.Error(
			"Failed to get meeting from meeting service: " +
				err.Error(),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get meeting",
		})
		return
	}

	participants := make([]entity.User, 0, len(resp.Participants))

	for _, p := range resp.Participants {
		participants = append(participants, entity.User{
			ID:   p.Id,
			Name: p.Name,
		})
	}

	c.JSON(http.StatusOK, entity.MeetingFull{
		ID:          resp.Id,
		Title:       resp.Title,
		Description: resp.Description,

		Creator: entity.User{
			ID:   resp.Creator.Id,
			Name: resp.Creator.Name,
		},

		Participants: participants,
		Address:       resp.Address,
		Latitude:      resp.Latitude,
		Longitude:     resp.Longitude,
		Type:          resp.Type,
		MeetingTime:   resp.MeetingTime,
		CurrentPeople: resp.CurrentPeople,
		MaxPeople:     resp.MaxPeople,
	})
}

func (h *Handler) ApplyOnMeeting(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        logger.Error("Invalid meeting id: " + err.Error())
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting id"})
        return
    }

    logger.Info("Applying on meeting: " + strconv.FormatInt(id, 10))

    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
        logger.Error("missing authorization header")
        return
    }

    userID, name, err := usecase.GetUserFromToken(c)
    if err != nil {
        logger.Error("Failed to get user from token: " + err.Error())
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
        return
    }

    _, err = h.meetingClient.ApplyOnMeeting(
        c.Request.Context(),
        &meeting.ApplyOnMeetingRequest{
            User: &meeting.User { 
                Id:   userID,
                Name: name,
            },
            MeetingId: id,
        },
    )

    if err != nil {
        logger.Error("Failed to apply on meeting: " + err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to apply on meeting"})
        return
    }

	_, err = h.chatClient.AddToMeetingChat(
        c.Request.Context(),
        &chat.CreateMeetingChatRequest{
            // TODO rename on user
            CreatorId:   userID,
        	CreatorName: name,
            MeetingId: id,
        },
    )

	if err != nil {
        logger.Error("Failed to add user to a meeting chat: " + err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to apply on meeting"})
        return
    }

    c.JSON(200, http.StatusOK)
}