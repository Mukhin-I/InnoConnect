package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.pg.innopolis.university/innoconnect-team/innoconnect/internal/gateway/entity"
	"gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/logger"
	"gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/pb/meeting"
)

func (h *Handler) CreateMeeting(c *gin.Context) {
	var req entity.CreateMeetingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Invalid JSON": err.Error(),
		})
		logger.Error("Invalid JSON in meeting request")
		return
	}

	meeting, err := h.meetingClient.CreateMeeting(
		c.Request.Context(),
		&meeting.CreateMeetingRequest{
			CreatorId:   1, // temporary
			CreatorName: "Pavel", // temporary
			Title:       req.Title,
			Description: req.Description,
			Type:        req.Type,
			MeetingTime: req.MeetingTime,
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