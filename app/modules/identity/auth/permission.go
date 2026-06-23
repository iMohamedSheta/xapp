package auth

import (
	"context"
	"fmt"

	"strings"
	"time"

	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xapp/app/x"
)

type PermissionService struct {
	permissionRepo *PermissionRepository
}

func NewPermissionService(permissionRepo *PermissionRepository) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
	}
}

func (s *PermissionService) getCacheKey(userID int64) string {
	return fmt.Sprintf("users:settings:permissions:%d", userID)
}

// GetGlobalPermissions retrieves all permissions for a user, including base role permissions
func (s *PermissionService) GetGlobalPermissions(ctx context.Context, userID int64, role enums.UserRole) (map[string]bool, error) {
	cacheKey := s.getCacheKey(userID)

	getPermissions := func() (map[string]bool, error) {
		// Load user-specific permissions from database
		userPermissions, err := s.permissionRepo.GetUserPermissionsMap(ctx, userID)
		if err != nil {
			// If there's an error loading permissions, return base permissions
			// (we don't cache the error, but we return base permissions for this request)
			return s.getBaseRolePermissions(role), nil
		}

		// If no specific permissions in DB, use base role permissions
		if len(userPermissions) == 0 {
			return s.getBaseRolePermissions(role), nil
		}

		return userPermissions, nil
	}

	// these roles always have the role base permissions
	if role == enums.RoleSuperAdmin || role == enums.RoleSuperManager {
		return s.getBaseRolePermissions(role), nil
	}

	return utils.CacheRemember(ctx, cacheKey, time.Minute, getPermissions)
}

// SavePermissions saves/updates permissions for a given user and clears the cache
func (s *PermissionService) SavePermissions(ctx context.Context, userID int64, permissions []string) error {
	if err := s.permissionRepo.SavePermissions(ctx, userID, permissions); err != nil {
		return err
	}

	// Invalidate cache
	_ = utils.CacheForget(ctx, s.getCacheKey(userID))
	return nil
}

// ClearUserPermissionCache clears the permission cache for a given user
func (s *PermissionService) ClearUserPermissionCache(ctx context.Context, userID int64) {
	_ = utils.CacheForget(ctx, s.getCacheKey(userID))
}

// ResetPermissions resets permissions for a given user and clears the cache
func (s *PermissionService) ResetPermissions(ctx context.Context, userID int64) error {
	if err := s.permissionRepo.ResetPermissions(ctx, userID); err != nil {
		return err
	}

	// Invalidate cache
	s.ClearUserPermissionCache(ctx, userID)
	return nil
}

// getBaseRolePermissions returns the base permissions for each role
// Loads permissions from config, with fallback to hardcoded defaults
func (s *PermissionService) getBaseRolePermissions(role enums.UserRole) map[string]bool {
	permissions := make(map[string]bool)

	// Try to load from config first
	configKey := "permission.default_permissions." + string(role)
	rawPermissions, _ := x.Config().Get(configKey)

	if rawPermissions != nil {
		// Flatten the result regardless of whether it's a map or a slice
		flatList := s.flattenPermissions(rawPermissions)
		for _, perm := range flatList {
			permissions[perm] = true
		}
	}

	return permissions
}

// flattenPermissions recursively converts nested permission maps/slices into a flat list of dot-notation strings.
// e.g., map[string]any{"offers": []any{"view", "create"}} -> []string{"offers.view", "offers.create"}
func (s *PermissionService) flattenPermissions(input any) []string {
	var result []string

	switch v := input.(type) {
	case string:
		result = append(result, v)
	case []string:
		result = append(result, v...)
	case []any:
		for _, item := range v {
			result = append(result, s.flattenPermissions(item)...)
		}
	case map[string]any:
		for category, actions := range v {
			flatActions := s.flattenPermissions(actions)
			for _, action := range flatActions {
				result = append(result, category+"."+action)
			}
		}
	}

	return result
}

// HasPermission checks if a user has a specific permission
func (s *PermissionService) HasPermission(ctx context.Context, userID int64, role enums.UserRole, permissionKey string) (bool, error) {
	permissions, err := s.GetGlobalPermissions(ctx, userID, role)
	if err != nil {
		return false, err
	}

	// Check for wildcard permission (super admin)
	if permissions["*"] {
		return true, nil
	}

	// Check exact permission
	if permissions[permissionKey] {
		return true, nil
	}

	// Check wildcard permissions (e.g., "admin.*" matches "admin.users.create")
	// Split permission key by dots and check progressively
	parts := splitPermissionKey(permissionKey)
	for i := len(parts); i > 0; i-- {
		wildcardKey := joinPermissionKey(parts[:i]) + ".*"
		if permissions[wildcardKey] {
			return true, nil
		}
	}

	return false, nil
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
