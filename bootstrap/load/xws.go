package load

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xioc"
	"github.com/imohamedsheta/xws"
)

func InitWebsocketServer(c *xioc.Container) {
	err := xioc.Singleton(c, func(c *xioc.Container) (*xws.Server, error) {
		originRaw := x.Config().GetWithDefault("cors.origin", nil)
		var origin []string
		if originRaw != nil {
			if val, ok := originRaw.([]string); ok {
				origin = val
			}
		}

		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return checkOrigin(r, origin) },
			// HandshakeTimeout: 1 * time.Second,
		}

		hub := xws.NewHub()
		adapter := xws.NewPusherAdapter(hub, "app-key", "app-secret")

		server := xws.NewWsServer(hub, upgrader, adapter, func(r *http.Request) (string, bool) {
			id := r.Header.Get("X-User-ID") // replace with JWT validation
			return id, id != ""
		})

		return server, nil
	})

	if err != nil {
		errMsg := "Failed to load websocket server module in the ioc container : " + err.Error()
		x.Logger().Error(errMsg)
	}
}

func checkOrigin(r *http.Request, allowedOrigins []string) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		// Some clients (like localhost testing) may not send Origin header
		return true
	}

	for _, o := range allowedOrigins {
		if o == "*" {
			return true
		}

		// Exact match
		if origin == o {
			return true
		}

		// Allow subdomains (like *.example.com)
		if len(o) > 2 && strings.HasPrefix(o, "*.") && strings.HasSuffix(origin, o[1:]) {
			return true
		}
	}

	// Fallback: allow private-network (LAN / loopback) origins
	host := utils.OriginHost(origin)
	if utils.IsPrivateOrLoopback(host) {
		return true
	}

	return false
}
