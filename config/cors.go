package config

import (
	"fmt"

	"github.com/imohamedsheta/xfig"
)

func init() {
	xfig.Register(corsConfig)
}
func corsConfig(cfg *xfig.Config) {
	appBindPort := xfig.Env("APP_BIND_PORT", "8080")
	cfg.Set("cors", map[string]any{
		"origin": []string{
			fmt.Sprintf("localhost:%s", appBindPort),
			fmt.Sprintf("127.0.0.1:%s", appBindPort),
		},
		"methods": []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		"allowed_headers": []string{
			"*",
		},
		"exposed_headers": []string{
			"X-XSRF-TOKEN",
		},
		"credentials":            true, // Allow cookies, HTTP auth, etc.
		"allow_private_networks": true,
		"max_age":                "86400", // Preflight cache in seconds
	})
}
