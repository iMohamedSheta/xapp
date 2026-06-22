//go:build dev

package load

import (
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/pkg/logger"
)

func buildLoggerManager() (*logger.Manager, error) {
	cfg := x.Config()
	defaultChannel := cfg.GetString("log.default", "app_log")
	channels := cfg.GetMap("log.channels", nil) // loads everything

	manager := logger.NewManager()
	var defaultLoaded bool

	for name, channelConfigRaw := range channels {
		if channelConfig, ok := channelConfigRaw.(map[string]any); ok {
			zapCfg := buildZapConfig(channelConfig)
			path := channelConfig["path"].(string)

			if name == defaultChannel {
				manager.LoadDefault(path, zapCfg)
				defaultLoaded = true
			} else {
				manager.Register(name, path, zapCfg)
			}
		}
	}

	ensureDefaultLoaded(manager, defaultLoaded, defaultChannel, channels)
	return manager, nil
}
