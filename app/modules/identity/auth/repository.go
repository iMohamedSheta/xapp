package auth

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/modules/settings"
	"github.com/imohamedsheta/xapp/app/shared/enums"
)

// PermissionRepository handles all data access operations for user permissions.
// Even though permissions are currently stored as JSON in the 'settings' table,
// having a dedicated repository provides a clean abstraction layer. This makes the
// codebase more maintainable and allows for easier changes if the storage mechanism
// (e.g., moving back to a normalized table or using a different store) changes in the future.
// This choice was by me @iMohamedSheta
type PermissionRepository struct {
	settingRepo *settings.SettingRepository
}

func NewPermissionRepository(settingRepo *settings.SettingRepository) *PermissionRepository {
	return &PermissionRepository{
		settingRepo: settingRepo,
	}
}

// SavePermissions saves/updates permissions for a given user in the settings table.
func (r *PermissionRepository) SavePermissions(c context.Context, userID int64, permissions []string) error {
	return r.settingRepo.SaveSetting(c, "user", &userID, enums.SettingPermissions, map[string]any{
		"permissions": permissions,
	}, nil)
}

// GetUserPermissions retrieves all permission keys for a given user
func (r *PermissionRepository) GetUserPermissions(c context.Context, userID int64) ([]string, error) {
	result, err := xqb.Table("settings").
		WithContext(c).
		Where("model", "=", "user").
		Where("model_id", "=", userID).
		Where("type", "=", "permissions").
		First()

	if err != nil {
		if errors.Is(err, xqb.ErrNotFound) {
			return []string{}, nil
		}
		return nil, err
	}

	var settings map[string]any
	rawSettings := result["settings"]

	if b, ok := rawSettings.([]byte); ok {
		_ = json.Unmarshal(b, &settings)
	} else if s, ok := rawSettings.(string); ok && s != "" {
		_ = json.Unmarshal([]byte(s), &settings)
	} else if m, ok := rawSettings.(map[string]any); ok {
		settings = m
	}

	if settings != nil {
		if perms, ok := settings["permissions"].([]any); ok {
			finalPerms := make([]string, 0, len(perms))
			for _, p := range perms {
				if s, ok := p.(string); ok {
					finalPerms = append(finalPerms, s)
				}
			}
			return finalPerms, nil
		}
	}

	return []string{}, nil
}

// GetUserPermissionsMap retrieves all permissions for a user as a map[string]bool for quick lookup
func (r *PermissionRepository) GetUserPermissionsMap(c context.Context, userID int64) (map[string]bool, error) {
	permissions, err := r.GetUserPermissions(c, userID)
	if err != nil {
		return nil, err
	}

	permissionsMap := make(map[string]bool)
	for _, perm := range permissions {
		permissionsMap[perm] = true
	}

	return permissionsMap, nil
}

// ResetPermissions removes custom permissions for a given user from the settings table.
// This effectively resets them to the role-based defaults.
func (r *PermissionRepository) ResetPermissions(c context.Context, userID int64) error {
	_, err := xqb.Table("settings").
		WithContext(c).
		Where("model", "=", "user").
		Where("model_id", "=", userID).
		Where("type", "=", enums.SettingPermissions).
		Delete()

	return err
}
