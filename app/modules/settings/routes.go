package settings

import (
	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/app/http/middleware"
	"github.com/imohamedsheta/xapp/app/x"
)

// AuthMiddleware is the minimal interface required for route registration.
type AuthMiddleware interface {
	Auth() gin.HandlerFunc
	SuperAdminOrAdminOnly() gin.HandlerFunc
	ClientOnly() gin.HandlerFunc
	ManagerOnly() gin.HandlerFunc
}

// RegisterRoutes registers settings routes onto the given router group.
func RegisterRoutes(r *gin.RouterGroup, authMiddleware AuthMiddleware) {
	authGroup := r.Group("/")
	authGroup.Use(authMiddleware.Auth())

	if settingsHandler, ok := x.App[*SettingsHandler](); ok && settingsHandler != nil {
		settingsGroup := authGroup.Group("/settings")
		settingsGroup.GET("/appearance", wI(settingsHandler.AppearanceView))
		settingsGroup.POST("/save", wI(settingsHandler.SaveSettings))
	}
}

// wI [WrapInertia] is a helper for wrapping error_handler and Inertia middleware together
// Fixed wI helper function
func wI(h func(*gin.Context) error) gin.HandlerFunc {
	return middleware.InertiaMiddlewareWithErrorHandler(h)
}

func e(h func(*gin.Context) error) gin.HandlerFunc {
	return middleware.HandleErrors(h)
}
