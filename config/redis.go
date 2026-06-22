package config

import "github.com/imohamedsheta/xfig"

func init() {
	xfig.Register(redisConfig)
}

// Redis configuration for the application
func redisConfig(cfg *xfig.Config) {
	cfg.Set("redis", map[string]any{
		// Default Redis connection to use
		"default": xfig.Env("REDIS_DEFAULT", "default"),

		// Global Redis options
		"options": map[string]any{
			"cluster": xfig.Env("REDIS_CLUSTER", "redis"),
			"prefix":  xfig.Env("REDIS_PREFIX", nil),
		},

		// Redis connections
		"connections": map[string]any{
			"default": map[string]any{
				// If there is url provided, it will be used instead of the configuration
				// and you can set all the other options in the url and include new options if they are valid
				"url":       xfig.Env("REDIS_URL", nil),
				"host":      xfig.Env("REDIS_HOST", "127.0.0.1"),
				"password":  xfig.Env("REDIS_PASSWORD", nil),
				"port":      xfig.Env("REDIS_PORT", 6379),
				"database":  xfig.Env("REDIS_DB", 10),
				"active":    xfig.Env("REDIS_DEFAULT_ACTIVE", true),
				"pool_size": xfig.Env("REDIS_POOL_SIZE", 0),
				"timeout":   xfig.Env("REDIS_TIMEOUT", ""),
			},
			"queue": map[string]any{
				"url":       xfig.Env("REDIS_QUEUE_URL", nil),
				"host":      xfig.Env("REDIS_QUEUE_HOST", "127.0.0.1"),
				"password":  xfig.Env("REDIS_QUEUE_PASSWORD", nil),
				"port":      xfig.Env("REDIS_QUEUE_PORT", 6379),
				"database":  xfig.Env("REDIS_QUEUE_DB", 9),
				"active":    xfig.Env("REDIS_QUEUE_ACTIVE", true),
				"pool_size": xfig.Env("REDIS_QUEUE_POOL_SIZE", 0),
				"timeout":   xfig.Env("REDIS_QUEUE_TIMEOUT", ""),
			},
		},
	})
}
