package config

import (
	"time"

	"github.com/imohamedsheta/xfig"
)

func init() {
	xfig.Register(appConfig)
}
func appConfig(cfg *xfig.Config) {
	cfg.Set("app", map[string]any{
		"name":              xfig.Env("APP_NAME", "GoCrudRestApi"),
		"url":               xfig.Env("APP_URL", "localhost"),
		"port":              xfig.Env("APP_PORT", "8080"),
		"bind_address":      xfig.Env("APP_BIND_ADDRESS", "0.0.0.0"),
		"bind_port":         xfig.Env("APP_BIND_PORT", "8080"),
		"shutdown_timeout":  20 * time.Second,
		"api_prefix":        xfig.Env("APP_API_PREFIX", "/api/v1"),
		"env":               xfig.Env("APP_ENV", "dev"),
		"debug":             xfig.Env("APP_DEBUG", false),
		"demo":              xfig.Env("DEMO_MODE", false),
		"secret":            xfig.Env("APP_SECRET", "hxdCTfhtkyJBVEdasdsadascxzsTmr401QqGl1111"),
		"global_rate_limit": xfig.Env("APP_GLOBAL_RATE_LIMIT", "100-M"), // 100 requests per minute
		"super_managers":    []map[string]any{
			// {
			// 	"email":    "admin@admin.com",
			// 	"password": "123456789",
			// 	"name":     "Mohamed Sheta",
			// },
		},
		// Add HTTPS settings
		"https": map[string]any{
			"enabled":  xfig.Env("APP_HTTPS_ENABLED", false),
			"certFile": xfig.Env("APP_HTTPS_CERT", "certs/server_local.crt"),
			"keyFile":  xfig.Env("APP_HTTPS_KEY", "certs/server_local.key"),
			"port":     xfig.Env("APP_HTTPS_PORT", "8080"),
		},
	})
}
