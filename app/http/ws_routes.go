package apphttp

import (
	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/app/x"
)

// RegisterWebSocketRoutes registers only WebSocket-specific routes.
// All channels (notifications, topology tests, etc.) share a single /ws/connect endpoint.
// Channel access is controlled by policies registered in app/registers.go.
func RegisterWebSocketRoutes(router *gin.RouterGroup, auth gin.HandlerFunc) {
	wsAdmin := router.Group("/")
	wsAdmin.Use(auth)
	{
		// Single WebSocket entry point — clients subscribe to channels dynamically.
		wsAdmin.GET("/ws/connect", gin.WrapF(x.WS().ServeWS))
	}
}
