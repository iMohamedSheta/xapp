package load

import (
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/pkg/logger"
	"github.com/imohamedsheta/xioc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(c *xioc.Container) {
	err := xioc.Singleton(c, func(c *xioc.Container) (*logger.Manager, error) {
		return buildLoggerManager()
	})
	if err != nil {
		x.Logger().Error("Failed to load logger module: " + err.Error())
	}
}

func buildZapConfig(channel map[string]any) zap.Config {
	levelStr, _ := channel["level"].(string)
	zapLevel := toZapLevel(levelStr)

	path, _ := channel["path"].(string)
	if path == "" {
		path = "storage/logs/app.json"
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      levelStr == "debug",
		Encoding:         "json",
		OutputPaths:      []string{path},
		ErrorOutputPaths: []string{path},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	return cfg
}

func toZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	case "panic":
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}

func ensureDefaultLoaded(m *logger.Manager, loaded bool, defaultName string, channels map[string]any) {
	if loaded {
		return
	}

	// set first channel as default channel if there is no default channel
	for name := range channels {
		// set first channel as default then stop
		_ = m.SetDefaultLogger(name)
		break
	}
}
