package config

import "github.com/imohamedsheta/xfig"

func init() {
	xfig.Register(logConfig)
}

type LogChannel string

const (
	AppLog      LogChannel = "app_log"
	RequestLog  LogChannel = "request_log"
	QueueLog    LogChannel = "queue_log"
	EventBusLog LogChannel = "bus_log"
)

// logConfig sets the logging configuration for the application.
func logConfig(cfg *xfig.Config) {
	cfg.Set("log", map[string]any{
		"default": string(AppLog),

		"channels": map[string]any{
			string(AppLog): map[string]any{
				"driver":   "daily",
				"path":     xfig.Env("APP_LOG_PATH", "storage/logs/log.json"),
				"level":    "debug",
				"max_size": 100,
				"max_age":  30, // Days
				"backup":   false,
			},
			string(RequestLog): map[string]any{
				"driver":   "daily",
				"path":     xfig.Env("APP_REQUEST_LOG_PATH", "storage/logs/request.json"),
				"level":    "debug",
				"max_size": 100,
				"max_age":  30,
				"backup":   false,
			},
			string(EventBusLog): map[string]any{
				"driver":   "daily",
				"path":     xfig.Env("APP_BUS_LOG_PATH", "storage/logs/bus.json"),
				"level":    "debug",
				"max_size": 100,
				"max_age":  30,
				"backup":   false,
			},
			string(QueueLog): map[string]any{
				"driver":   "daily",
				"path":     xfig.Env("APP_QUEUE_LOG_PATH", "storage/logs/queue.json"),
				"level":    "debug",
				"max_size": 100,
				"max_age":  30,
				"backup":   false,
			},
		},
	})
}
