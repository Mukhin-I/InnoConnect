package app

import (
	"innoconnect/internal/gateway/grpcclient"
	"innoconnect/internal/gateway/transport"
	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "innoconnect/docs"
)

// Setup and run gateway server
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

	chatClient, err := grpcclient.NewChatClient()
	userClient, err := grpcclient.NewUserClient()
	if err != nil {
		logger.Error("Server startup failed" + err.Error())
		return err
	}

	wsHub := transport.NewHub()
	handler := transport.NewHandler(meetingClient, requestClient, chatClient, userClient, wsHub)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://10.93.27.21:5173",
			"http://10.93.27.21",
		},
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

	for _, route := range router.Routes() {
		logger.Info(route.Method + " " + route.Path)
	}

	router.Run(":" + gatewayPort)

	return nil
}

// Configure all HTTP endpoints
func setEndpoints(router *gin.Engine, h *transport.Handler) {
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
	router.GET("/me", h.GetCurrentUser)

	router.POST("/requests", h.CreateRequest)
	router.GET("/requests", h.GetRequests)
	router.GET("/requests/:id", h.GetRequest)
	router.DELETE("/requests/:id", h.DeleteRequest)
	router.POST("/requests/:id/chat", h.GetOrCreateRequestChat)
	router.POST("/requests/:id/apply", h.ApplyToRequest)

	router.DELETE(
		"/requests/:request_id/applications/:user_id",
		h.CancelRequestApplication,
	)

	router.GET("/meetings/:id/chat", h.GetMeetingChat)

	router.GET("/chats", h.GetChats)
	router.GET("/chats/:chat_id", h.GetChat)
	router.GET("/chats/:chat_id/messages", h.GetMessages)
	router.POST("/chats/:chat_id/messages", h.SendMessage)

	router.POST("/meetings", h.CreateMeeting)
	router.GET("/meetings", h.GetMeetings)
	router.GET("/meetings/:id", h.GetMeeting)
	router.DELETE("/meetings/:id", h.DeleteMeeting)
	router.POST("/meetings/:id", h.ApplyOnMeeting)

	router.GET("/ws", h.WebSocket)

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
}
