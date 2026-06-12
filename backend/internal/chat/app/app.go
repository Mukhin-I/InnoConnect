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
	chatPort := config.GetVar("CHAT_SERVICE_PORT")

	logger.Info("Char service running on :" + chatPort)

	router.Run(":" + chatPort)
}

func setEndpoints(router *gin.Engine) {
	// Here should be endpoints. For example:
	// router.GET("/ws", transport.WsHandler)
}
