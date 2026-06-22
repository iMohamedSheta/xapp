package orders

import "github.com/gin-gonic/gin"

// AuthMiddleware is the minimal interface required for route registration.
type AuthMiddleware interface {
	Auth() gin.HandlerFunc
	SuperAdminOrAdminOnly() gin.HandlerFunc
	ClientOnly() gin.HandlerFunc
	ManagerOnly() gin.HandlerFunc
}

// RegisterRoutes registers orders routes onto the given router group.
func RegisterRoutes(r *gin.RouterGroup, auth AuthMiddleware) {
	// TODO: add orders routes here
}
