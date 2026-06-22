package config

import (
	"time"

	"github.com/imohamedsheta/xfig"
)

func init() {
	xfig.Register(authConfig)
}

// authConfig sets the authentication configuration for the application.
func authConfig(cfg *xfig.Config) {
	cfg.Set("auth", map[string]any{
		"default": "jwt",
		"jwt": map[string]any{
			"secret":    "secret",
			"issuer":    xfig.Env("APP_NAME", "xapp"),
			"audience":  xfig.Env("APP_NAME", "xapp"),
			"algorithm": "HS256",
			"access_token": map[string]any{
				"expiry": 168 * time.Hour,
				"path":   "/",
			},
			"refresh_token": map[string]any{
				"expiry": 168 * time.Hour, // 7 days
				"path":   "/auth/refresh",
			},
		},
		"xsocial": map[string]any{
			"facebook": map[string]any{
				"client_id":     xfig.Env("FACEBOOK_OAUTH_CLIENT_ID", ""),
				"client_secret": xfig.Env("FACEBOOK_OAUTH_CLIENT_SECRET", ""),
				"redirect":      "http://localhost:8080/auth/callback/facebook",
				"scopes":        nil,
			},
			"google": map[string]any{
				"client_id":     xfig.Env("GOOGLE_OAUTH_CLIENT_ID", ""),
				"client_secret": xfig.Env("GOOGLE_OAUTH_CLIENT_SECRET", ""),
				"redirect":      "",
				"scopes":        nil,
			},
			"github": map[string]any{
				"client_id":     xfig.Env("GITHUB_OAUTH_CLIENT_ID", ""),
				"client_secret": xfig.Env("GITHUB_OAUTH_CLIENT_SECRET", ""),
				"redirect":      "http://localhost:8080/auth/callback/github",
				"scopes":        nil,
			},
		},
	})
}
