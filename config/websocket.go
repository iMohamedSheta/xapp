package config

import (
	"time"

	"github.com/imohamedsheta/xfig"
)

func init() {
	xfig.Register(websocketConfig)
}

func websocketConfig(cfg *xfig.Config) {
	cfg.Set("websocket", map[string]any{
		"name":             xfig.Env("APP_NAME", "xapp"),
		"url":              xfig.Env("APP_URL", "localhost"),
		"port":             xfig.Env("APP_PORT", "8080"),
		"bind_address":     xfig.Env("APP_BIND_ADDRESS", "0.0.0.0"),
		"bind_port":        xfig.Env("APP_WEBSOCKET_BIND_PORT", "8081"),
		"shutdown_timeout": 20 * time.Second,
		"env":              xfig.Env("APP_ENV", "dev"),
		"debug":            xfig.Env("APP_DEBUG", false),
	})
}
