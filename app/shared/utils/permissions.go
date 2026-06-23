package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/shared/enums"
)

// GetPermissions retrieves the permissions map from the context
func GetPermissions(c *gin.Context) (map[string]bool, error) {
	permissionsData, ok := c.Get(string(enums.ContextKeyPermissions))
	if !ok {
		return make(map[string]bool), nil
	}

	permissions, ok := permissionsData.(map[string]bool)
	if !ok {
		return nil, xerr.New("failed to get permissions from context", enums.XErrServerError, nil)
	}

	return permissions, nil
}

// HasPermission checks if the user has a specific permission
func HasPermission(c *gin.Context, permissionKey string) (bool, error) {
	permissions, err := GetPermissions(c)
	if err != nil {
		return false, err
	}

	// Check for wildcard permission
	if permissions["*"] {
		return true, nil
	}

	// Check exact permission
	if permissions[permissionKey] {
		return true, nil
	}

	// Check wildcard permissions (e.g., "admin.*" matches "admin.users.create")
	parts := splitPermissionKey(permissionKey)
	for i := len(parts); i > 0; i-- {
		wildcardKey := joinPermissionKey(parts[:i]) + ".*"
		if permissions[wildcardKey] {
			return true, nil
		}
	}

	return false, nil
}

// MustHavePermission checks if the user has a specific permission and returns an error if not
func MustHavePermission(c *gin.Context, permissionKey string) error {
	hasPermission, err := HasPermission(c, permissionKey)
	if err != nil {
		return err
	}

	if !hasPermission {
		return xerr.New("insufficient permissions", enums.XErrForbiddenError, nil)
	}

	return nil
}

// Helper function to split permission key by dots
func splitPermissionKey(key string) []string {
	parts := make([]string, 0)
	current := ""
	for _, char := range key {
		if char == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// Helper function to join permission key parts
func joinPermissionKey(parts []string) string {
	var result strings.Builder
	for i, part := range parts {
		if i > 0 {
			result.WriteString(".")
		}
		result.WriteString(part)
	}
	return result.String()
}
