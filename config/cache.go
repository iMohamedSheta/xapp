package config

import "github.com/imohamedsheta/xfig"

func init() {
	xfig.Register(cacheConfig)
}
func cacheConfig(cfg *xfig.Config) {
	cfg.Set("cache", map[string]any{
		"default": "redis",
		"stores": map[string]any{
			"redis": map[string]any{
				"driver": "redis",
				"host":   "127.0.0.1",
				"port":   "6379",
			},
		},
	})
}
