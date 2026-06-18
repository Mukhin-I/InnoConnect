package transport

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"innoconnect/internal/gateway/entity"
	"innoconnect/pkg/logger"
	"innoconnect/pkg/pb/meeting"
)

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

	meeting, err := h.meetingClient.CreateMeeting(
		c.Request.Context(),
		&meeting.CreateMeetingRequest{
			CreatorId:   1, // temporary
			CreatorName: "Pavel", // temporary
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

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logger.Error("Falied go get answer from meeting service via gRPC" + err.Error())
		return
	}

	c.JSON(http.StatusCreated, meeting)
}

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