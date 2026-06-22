package config

import (
	"github.com/imohamedsheta/xfig"
)

func init() {
	xfig.Register(databaseConfig)
}

func databaseConfig(cfg *xfig.Config) {
	cfg.Set("database", map[string]any{

		// This is the default database connection should be valid connection to use.
		"default": xfig.Env("DB_CONNECTION", "default"),

		"connections": map[string]any{
			// Postgres connection
			"default": map[string]any{
				"host":     xfig.Env("DB_HOST", "localhost"),
				"port":     xfig.Env("DB_PORT", "5432"),
				"user":     xfig.Env("DB_USERNAME", "postgres"),
				"pass":     xfig.Env("DB_PASSWORD", "123456"),
				"database": xfig.Env("DB_DATABASE", "xapp"),
				"dialect":  "postgres",
				"sslmode":  "disable",
				"timezone": "UTC",
			},
		},
		// connection pool settings
		"max_idle_conns":    xfig.Env("DB_MAX_IDLE_CONNS", 10),
		"max_open_conns":    xfig.Env("DB_MAX_OPEN_CONNS", 100),
		"conn_max_lifetime": xfig.Env("DB_CONN_MAX_LIFETIME", "30m"),
	})
}
