package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xapp/app/x"
)

// User returns the authenticated user from the gin context.
// Delegates to x.AuthUser — use x.AuthUser directly to avoid importing this package.
func User(c *gin.Context) (*models.User, *xerr.XErr) {
	return x.AuthUser(c)
}

// ClearAuthCookies clears the authentication cookies.
func ClearAuthCookies(c *gin.Context) {
	x.ClearAuthCookies(c)
}
