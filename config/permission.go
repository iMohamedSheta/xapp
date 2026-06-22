package config

import (
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xfig"
)

func init() {
	xfig.Register(permissionConfig)
}

func permissionConfig(cfg *xfig.Config) {
	cfg.Set("permission", map[string]any{
		"default_permissions": map[string]any{
			string(enums.RoleSuperAdmin): map[string]any{
				"dashboard": []string{
					"view",
					"stats",
				},
			},
			string(enums.RoleAdmin): map[string]any{
				"dashboard": []string{
					"view",
					"stats",
				},
			},
			string(enums.RoleSuperManager): map[string]any{
				"dashboard": []string{"view"},
				"plans":     []string{"view", "create", "update", "delete"},
				"tenants":   []string{"view", "create", "update", "delete"},
				"users":     []string{"view", "create", "update", "delete"},
				"invoices":  []string{"view", "pay"},
			},
			string(enums.RoleManager): map[string]any{
				"dashboard": []string{"view"},
				"plans":     []string{"view"},
				"tenants":   []string{"view"},
				"users":     []string{"view"},
				"invoices":  []string{"view"},
			},
		},
		"available_permissions": map[string]any{
			"admin_group": map[string]any{ // super admin (have full permissions), admin (have full permissions as default)
				"dashboard": []string{
					"view",
					"stats",
				},
			},
		},
	})
}
