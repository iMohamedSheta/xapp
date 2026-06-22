package app

import (
	"maps"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	apphttp "github.com/imohamedsheta/xapp/app/http"
	"github.com/imohamedsheta/xapp/app/http/middleware"
	"github.com/imohamedsheta/xapp/app/modules/identity/auth"
	"github.com/imohamedsheta/xapp/app/providers"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/bootstrap"
	"github.com/imohamedsheta/xapp/bootstrap/adapters"
	"github.com/imohamedsheta/xapp/pkg/chainq"
	"github.com/imohamedsheta/xioc"
)

// WebsocketHooks builds the WebsocketHooks struct injected into bootstrap.RunWebsocket.
func WebsocketHooks() bootstrap.WebsocketHooks {
	return bootstrap.WebsocketHooks{
		RegisterRoutes:      apphttp.RegisterWebSocketRoutes,
		BuildAuthMiddleware: buildWebSocketAuthMiddleware,
		Middleware: bootstrap.GinMiddleware{
			Recovery: middleware.RecoveryWithLogger,
			Logger:   middleware.Logger,
			CORS:     middleware.CORSMiddleware,
		},
	}
}

// WebSocketChannels configures all WebSocket channel access policies.
// Call once after the WS server is booted.
func WebSocketChannels() {
	hub := x.WS().Hub

	// register all websocket channels
	for _, channel := range websocketChannels {
		hub.RegisterChannel(channel)
	}
}

// ServiceProviders registers all IoC container bindings.
func ServiceProviders(c *xioc.Container) {
	providers.RegisterProviders(c)
	// Compile-time check: panic early if any x.App[T]() binding is missing.
	xioc.MustAllRegistered(
		iocMustAllRegistered...,
	)
}

var (
	builtHandlers map[string]asynq.Handler
	buildOnce     sync.Once
)

// TaskHandlers returns the full handler map (individual + chain orchestrator).
// Memoized — safe to call from RunWorker and RunWebsocket.
func TaskHandlers() map[string]asynq.Handler {
	buildOnce.Do(func() {
		orchestrator := chainq.NewChainOrchestrator(
			x.Queue().Client,
			adapters.NewLoggerAdapter(x.Logger()),
		)

		individual := taskHandlers()
		for taskType, handler := range individual {
			orchestrator.RegisterHandler(taskType, handler)
		}

		builtHandlers = make(map[string]asynq.Handler)
		maps.Copy(builtHandlers, individual)
		builtHandlers[chainq.TypeChainOrchestrator] = orchestrator
	})

	return builtHandlers
}

func buildWebSocketAuthMiddleware() gin.HandlerFunc {
	return x.AppMust[*auth.AuthMiddleware]().WebSocketAuth()
}
