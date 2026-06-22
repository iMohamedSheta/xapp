package config

import "github.com/imohamedsheta/xfig"

func init() {
	xfig.Register(queueConfig)
}

// Redis configuration for the application
func queueConfig(cfg *xfig.Config) {
	cfg.Set("queue", map[string]any{
		// Default Redis connection to use, we will use redis connection (jobs) if it's redis
		"default":  xfig.Env("QUEUE_DEFAULT", "redis"),
		"enabled":  xfig.Env("QUEUE_ACTIVE", true),
		"required": xfig.Env("QUEUE_REQUIRED", false), // Required for app to run
		"consumer": map[string]any{
			// Worker concurrency settings
			"concurrency": xfig.Env("QUEUE_CONSUMER_CONCURRENCY", 10),

			// Queue priorities
			"queues": map[string]any{
				"critical":      xfig.Env("QUEUE_PRIORITY_CRITICAL", 6),
				"default":       xfig.Env("QUEUE_PRIORITY_DEFAULT", 3),
				"low":           xfig.Env("QUEUE_PRIORITY_LOW", 1),
				"radius":        6,
				"notifications": 3,
			},
			"websocket_queues": map[string]any{
				"ws_notifications": 2,
			},
			// Retry configuration
			"retry": map[string]any{
				"max_attempts": xfig.Env("QUEUE_MAX_RETRY_ATTEMPTS", 3),
				"delay":        xfig.Env("QUEUE_RETRY_DELAY", "15s"), // delay between retries
			},

			// Health check configuration
			"health_check": map[string]any{
				"enabled":  xfig.Env("QUEUE_HEALTH_CHECK_ENABLED", true),
				"interval": xfig.Env("QUEUE_HEALTH_CHECK_INTERVAL", "30s"),
			},

			// Logging configuration for tasks
			"logging": map[string]any{
				"log_level":        xfig.Env("QUEUE_LOG_LEVEL", "info"),
				"log_failed_tasks": xfig.Env("QUEUE_LOG_FAILED_TASKS", true),
				"log_success":      xfig.Env("QUEUE_LOG_SUCCESS", true),
			},
		},
	})
}
