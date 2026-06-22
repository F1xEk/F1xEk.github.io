package routes

import (
	"charging-app/handlers"
	"charging-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	api := router.Group("/api")

	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		api.GET("/ws", handlers.WSHandler)
	}

	protected := api.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())

	{
	
	protected.GET("/profile", handlers.GetProfile)
	protected.GET("/stations", handlers.GetStations)
	

	protected.POST(
		"/start-session",
		handlers.StartSession,
	)

	protected.GET(
		"/sessions",
		handlers.GetSessions,
	)

	protected.POST(
		"/stop-session",
		handlers.StopSession,
	)

	protected.GET(
	"/analytics",
	handlers.GetAnalytics,
	)
	}
}