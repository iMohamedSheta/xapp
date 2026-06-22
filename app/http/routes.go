package apphttp

import (
	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/http/middleware"
	"github.com/imohamedsheta/xapp/app/modules/identity"
	"github.com/imohamedsheta/xapp/app/modules/identity/auth"
	"github.com/imohamedsheta/xapp/app/modules/settings"
	"github.com/imohamedsheta/xapp/app/x"
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	r.NoRoute(e(func(ctx *gin.Context) error {
		return xerr.New("route does not exists", enums.XErrNotFoundError, nil)
	}))

	r.NoMethod(e(func(ctx *gin.Context) error {
		return xerr.New("method does not exists", enums.XErrNotFoundError, nil)
	}))

	// serve static assets
	// r.Static("/build", "./public/build")

	r.Static("/public", "./public")
	r.Static("/storage", "./storage/app/public")

	// Global middlewares
	r.Use(middleware.RecoveryWithLogger())
	r.Use(middleware.Logger())

	// Register Api routes here
	registerApiRoutes(r)

	// r.Use(middleware.RateLimiter())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.CSRF())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "Everything is working fine",
		})
	})

	// auto redirect to login in home page
	r.GET("/", utils.Redirect("/login"))

	// register modules routes
	registerModulesRoutes(r)

	return r
}

func registerModulesRoutes(r *gin.Engine) {
	authMiddleware := x.AppMust[*auth.AuthMiddleware]()

	// register identity routes
	identityGroup := r.Group("/")
	identity.RegisterRoutes(identityGroup, authMiddleware)

	// register settings routes
	settingsGroup := r.Group("/")
	settings.RegisterRoutes(settingsGroup, authMiddleware)
}

func registerApiRoutes(r *gin.Engine) {
	// register api routes here
}

// wI [WrapInertia] is a helper for wrapping error_handler and Inertia middleware together
// Fixed wI helper function
func wI(h func(*gin.Context) error) gin.HandlerFunc {
	return middleware.InertiaMiddlewareWithErrorHandler(h)
}

func e(h func(*gin.Context) error) gin.HandlerFunc {
	return middleware.HandleErrors(h)
}
