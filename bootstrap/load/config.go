package load

import (
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xioc"

	"github.com/imohamedsheta/xfig"
)

func InitConfig(c *xioc.Container) {
	err := xioc.Singleton(c, func(c *xioc.Container) (*xfig.Config, error) {
		cfg := xfig.New()                // Global config
		xfig.ApplyRegisteredLoaders(cfg) // Apply config module registered loaders
		return cfg, nil
	})

	if err != nil {
		x.Logger().Error("Failed to load config module as singleton in the ioc container: " + err.Error())
	}
}
