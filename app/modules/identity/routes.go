package identity

import (
	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/app/http/middleware"
	"github.com/imohamedsheta/xapp/app/modules/identity/auth"
	"github.com/imohamedsheta/xapp/app/modules/identity/users"
	"github.com/imohamedsheta/xapp/app/x"
)

// AuthMiddleware is the minimal interface required for route registration.
// This avoids a circular import with app/http/middleware.
type AuthMiddleware interface {
	Auth() gin.HandlerFunc
	RedirectToDashboardIfAuthenticated() gin.HandlerFunc
	SuperAdminOrAdminOnly() gin.HandlerFunc
	ClientOnly() gin.HandlerFunc
	ManagerOnly() gin.HandlerFunc
}

// RegisterRoutes registers auth routes onto the given router group.
func RegisterRoutes(r *gin.RouterGroup, authMiddleware AuthMiddleware) {
	web := r.Group("/")
	web.Use(authMiddleware.RedirectToDashboardIfAuthenticated())
	if authHandler, ok := x.App[*auth.AuthHandler](); ok && authHandler != nil {
		web.GET("/login", wI(authHandler.LoginView))
		web.GET("/register", wI(authHandler.RegisterView))

		web.POST("/login", wI(authHandler.Login))
		web.POST("/register", wI(authHandler.Register))
	}

	authGroup := r.Group("/")
	authGroup.Use(authMiddleware.Auth())

	if authHandler, ok := x.App[*auth.AuthHandler](); ok && authHandler != nil {
		authGroup.GET("/dashboard", wI(authHandler.DashboardView))
		authGroup.POST("/logout", wI(authHandler.Logout))
	}

	if accountHandler, ok := x.App[*users.AccountHandler](); ok && accountHandler != nil {
		settingsGroup := authGroup.Group("/settings")
		settingsGroup.GET("/profile", wI(accountHandler.ProfileView))
		settingsGroup.PATCH("/profile/update", wI(accountHandler.UpdateProfile))
		settingsGroup.GET("/password", wI(accountHandler.PasswordView))
		settingsGroup.PATCH("/password/update", wI(accountHandler.UpdatePassword))
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
