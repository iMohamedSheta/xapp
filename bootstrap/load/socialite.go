package load

import (
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xioc"
	"github.com/imohamedsheta/xsocial"
)

func InitSocialite(c *xioc.Container) {
	err := xioc.Singleton(c, func(c *xioc.Container) (*xsocial.Socialite, error) {
		s := xsocial.New()
		cfg := x.Config()

		socialMap := cfg.GetMap("auth.socialite", nil)
		for name, raw := range socialMap {
			if m, ok := raw.(map[string]any); ok {

				// scopes may be nil OR []string
				var scopes []string
				if rawScopes, found := m["scopes"]; found && rawScopes != nil {
					if tmp, ok := rawScopes.([]string); ok {
						scopes = tmp
					}
				}

				s.AddConfig(name, xsocial.Config{
					ClientID:     utils.ToString(m["client_id"]),
					ClientSecret: utils.ToString(m["client_secret"]),
					RedirectURL:  utils.ToString(m["redirect"]),
					Scopes:       scopes,
				})
			}
		}

		return s, nil
	})

	if err != nil {
		x.Logger().Error("Failed to load socialite module as singleton in the ioc container: " + err.Error())
	}
}
