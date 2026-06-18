package app

import (
	"github.com/gin-gonic/gin"
	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"
)

// Function for setuping Gin server
func CreateServer() {
	router := gin.Default()

	setEndpoints(router)

	// Getting .env vars
	userPort := config.GetVar("USER_SERVICE_PORT")

	logger.Info("User service running on :" + userPort)

	router.Run(":" + userPort)
}

func setEndpoints(router *gin.Engine) {
	// Here should be endpoints. For example:
	// router.GET("/ws", transport.WsHandler)
}
