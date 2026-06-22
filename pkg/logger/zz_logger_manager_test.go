package logger_test

import (
	"sync"
	"testing"

	"github.com/imohamedsheta/xapp/pkg/logger"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newTestZapConfig(path string) zap.Config {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	}
	if path != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, path)
	}
	return cfg
}

func TestRegisterAndChannel(t *testing.T) {
	logManager := logger.NewManager()
	err := logManager.Register("test", "", newTestZapConfig(""))
	assert.NoError(t, err)

	l := logManager.Channel("test")
	assert.NotNil(t, l)
	l.Debug("Test log message")
}

func TestLoadDefault_Once(t *testing.T) {
	logManager := logger.NewManager()

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	logger1 := logManager.LoadDefault("", cfg)
	logger2 := logManager.Default()

	assert.Equal(t, logger1, logger2)
}

func TestSetDefaultLogger(t *testing.T) {
	logManager := logger.NewManager()
	_ = logManager.Register("log1", "", newTestZapConfig(""))
	_ = logManager.Register("log2", "", newTestZapConfig(""))

	_ = logManager.SetDefaultLogger("log2")
	assert.Equal(t, logManager.Channel("log2"), logManager.Default())
}

func TestConcurrentRegisterAndGet(t *testing.T) {
	logManager := logger.NewManager()
	wg := sync.WaitGroup{}
	num := 100

	for i := range num {
		name := "log" + string(rune(i+'A'))

		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			err := logManager.Register(name, "", newTestZapConfig(""))
			assert.NoError(t, err)

			logger := logManager.Channel(name)
			assert.NotNil(t, logger)
		}(name)
	}

	wg.Wait()
	assert.Equal(t, num, len(logManager.GetLoggers()))
}
