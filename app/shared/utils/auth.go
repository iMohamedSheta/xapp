package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/shared/enums"
)

func AuthSessionId(c *gin.Context) (string, error) {
	sessionId := c.GetString(string(enums.ContextKeySessionId))
	if sessionId == "" {
		return "", xerr.New("session id is empty in context", enums.XErrUnAuthorizedError, nil)
	}

	return sessionId, nil
}

func IsSuperAdminOrAdmin(c *gin.Context) bool {
	role, ok := c.Get(string(enums.ContextKeyRole))
	if !ok {
		return false
	}

	r, ok := role.(enums.UserRole)
	if !ok {
		return false
	}

	return r == enums.RoleAdmin || r == enums.RoleSuperAdmin
}

// IsSuperAdmin checks if the user has the super admin role.
func IsSuperAdmin(c *gin.Context) bool {
	role, ok := c.Get(string(enums.ContextKeyRole))
	if !ok {
		return false
	}

	r, ok := role.(enums.UserRole)
	if !ok {
		return false
	}

	return r == enums.RoleSuperAdmin
}

func IsClient(c *gin.Context) bool {
	role, ok := c.Get(string(enums.ContextKeyRole))
	if !ok {
		return false
	}

	r, ok := role.(enums.UserRole)
	if !ok {
		return false
	}

	return r == enums.RoleClient
}

func IsManager(c *gin.Context) bool {
	role, ok := c.Get(string(enums.ContextKeyRole))
	if !ok {
		return false
	}
	r, ok := role.(enums.UserRole)
	if !ok {
		return false
	}
	return r == enums.RoleManager || r == enums.RoleSuperManager
}
