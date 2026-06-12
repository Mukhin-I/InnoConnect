package app

import (
	"github.com/gin-gonic/gin"
	"gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/config"
	"gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/logger"
)

// Function for setuping Gin server
func CreateServer() {
	router := gin.Default()

	setEndpoints(router)

	// Getting .env vars
	meetingPort := config.GetVar("MEETING_SERVICE_PORT")

	logger.Info("Meeting service running on :" + meetingPort)

	router.Run(":" + meetingPort)
}

func setEndpoints(router *gin.Engine) {
	// Here should be endpoints. For example:
	// router.GET("/ws", transport.WsHandler)
}
