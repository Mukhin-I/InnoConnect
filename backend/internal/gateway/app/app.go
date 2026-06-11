package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/logger"
)

// Function for setuping Gin server
func CreateServer() {
	router := gin.Default()

	setEndpoints(router)

	// Getting .env vars
	godotenv.Load()
	gatewayPort := os.Getenv("GATEWAY_SERVICE_PORT")

	logger.Info("Gateway service running on :" + gatewayPort)

	router.Run(":" + gatewayPort)
}

func setEndpoints(router *gin.Engine) {
	// Here should be endpoints. For example:
	// router.GET("/ws", transport.WsHandler)
}
