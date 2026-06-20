package app

import (
	"github.com/gin-gonic/gin"
	"innoconnect/internal/gateway/grpcclient"
	"innoconnect/internal/gateway/transport"
	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"
	"github.com/gin-contrib/cors"
)

// Function for setuping Gin server
func CreateServer() error {
	meetingClient, err := grpcclient.NewMeetingClient()
	if err != nil {
		logger.Error("Server startup failed" + err.Error())
		return err
	}

	requestClient, err := grpcclient.NewRequestClient()
	if err != nil {
		logger.Error("Server startup failed" + err.Error())
		return err
	}

	handler := transport.NewHandler(meetingClient, requestClient)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"}, // your frontend
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	setEndpoints(router, handler)

	// Getting .env vars
	gatewayPort := config.GetVar("GATEWAY_SERVICE_PORT")

	logger.Info("Gateway service running on :" + gatewayPort)

	router.Run(":" + gatewayPort)

	return nil
}

func setEndpoints(router *gin.Engine, h *transport.Handler) {
	router.POST("/meetings", h.CreateMeeting)
	router.GET("/meetings", h.GetMeetings)
	router.GET("/meetings/:id", h.GetMeeting)

	router.POST("/requests", h.CreateRequest)
	router.GET("/requests", h.GetRequests)
	router.GET("/requests/:id", h.GetRequest)
}
