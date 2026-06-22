package settings

import (
	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/app/domain/utils"
)

type SettingPolicy struct {
}

// NewSettingPolicy creates a new SettingPolicy with context
func NewSettingPolicy() *SettingPolicy {
	return &SettingPolicy{}
}

// Check if user can view app settings
func (p *SettingPolicy) CanView(c *gin.Context) bool {
	hasPermission, err := utils.HasPermission(c, "app_settings.view")
	if err != nil {
		return false
	}
	return hasPermission
}

// Check if user can update app settings
func (p *SettingPolicy) CanUpdate(c *gin.Context) bool {
	hasPermission, err := utils.HasPermission(c, "app_settings.update")
	if err != nil {
		return false
	}
	return hasPermission
}
