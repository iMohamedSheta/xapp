package x

import (
	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xapp/app/shared/enums"
)

// AuthUser returns the authenticated user from the gin context.
// Returns an error if the user is not set (unauthenticated request).
func AuthUser(c *gin.Context) (*models.User, *xerr.XErr) {
	userData, ok := c.Get(string(enums.ContextKeyAuthUser))
	if !ok {
		return nil, xerr.New("missing auth user data in context", enums.XErrUnAuthorizedError, nil)
	}

	user, ok := userData.(*models.User)
	if !ok {
		return nil, xerr.New("failed to bind authorized user data from context to user model", enums.XErrServerError, nil)
	}

	return user, nil
}

// ClearAuthCookies clears the authentication cookies from the response.
func ClearAuthCookies(c *gin.Context) {
	domain := ""
	c.SetCookie("access_token", "", -1, "/", domain, false, true)
	c.SetCookie("refresh_token", "", -1, "/auth/refresh", domain, false, true)
}
