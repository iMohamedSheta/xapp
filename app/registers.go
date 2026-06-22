package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"github.com/imohamedsheta/xapp/app/cmd"
	"github.com/imohamedsheta/xapp/app/domain/events"
	"github.com/imohamedsheta/xapp/app/domain/notifications/handlers"
	"github.com/imohamedsheta/xapp/app/domain/rules"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/modules/audit_logs"
	audit_logs_listeners "github.com/imohamedsheta/xapp/app/modules/audit_logs/listeners"
	"github.com/imohamedsheta/xapp/app/modules/identity/auth"
	"github.com/imohamedsheta/xapp/app/modules/identity/users"
	notifications_listeners "github.com/imohamedsheta/xapp/app/modules/notifications/listeners"
	"github.com/imohamedsheta/xapp/app/modules/settings"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/pkg/bus"
	asynqBus "github.com/imohamedsheta/xapp/pkg/bus/asynq_backned"
	"github.com/imohamedsheta/xioc"
	"github.com/imohamedsheta/xnotify"
	asynqNotify "github.com/imohamedsheta/xnotify/asynq"
	"github.com/imohamedsheta/xws"
	"github.com/spf13/cobra"
)

/*
Package registers is the single place to wire up the application.

Add new:
  - Validation rules   → ValidationRules()
  - Task handlers      → TaskHandlers()
  - Notify channels    → NotifyChannels()
  - WebSocket channels → WebSocketChannels()
  - Service providers  → ServiceProviders()
  - CLI commands       → Commands()
*/

// ValidationRules defines all custom validation rules.
func ValidationRules() map[string]validator.FuncCtx {
	return map[string]validator.FuncCtx{
		"unique_db":      rules.UniqueInDB,
		"exists_db":      rules.ExistsInDB,
		"egyptian_phone": rules.EgyptianPhone,
	}
}

// NotifyChannels defines all notification channel handlers.
func NotifyChannels() map[string]xnotify.ChannelHandler {
	return map[string]xnotify.ChannelHandler{
		"database":  handlers.DatabaseChannelHandler,
		"websocket": handlers.WebsocketChannelHandler,
		"whatsapp":  handlers.WhatsappChannelHandler,
	}
}

// Commands defines user defined CLI commands.
func Commands() []*cobra.Command {
	return []*cobra.Command{
		cmd.InspireCommand,
		cmd.SeedCommand,
	}
}

// websocketChannels configures all WebSocket channel access policies.
var websocketChannels = []*xws.ChannelPolicy{
	{
		Pattern: "user_notifications.*",
		CanRead: func(userID, channel string) bool {
			return channel == "user_notifications."+userID
		},
		CanWrite: func(userID, channel string) bool {
			return false // server-only writes via hub.Broadcast
		},
	},
}

// Compile time check: panic early if any x.App[T]() binding is missing to prevent runtime errors.
var iocMustAllRegistered = []string{
	xioc.TypeKey[*auth.JwtService](),
	xioc.TypeKey[*settings.SettingRepository](),
	xioc.TypeKey[*auth.PermissionRepository](),
	xioc.TypeKey[*auth.PermissionService](),
	xioc.TypeKey[*auth.AuthMiddleware](),
	xioc.TypeKey[*auth.AuthHandler](),
}

// BusListeners subscribes every domain event listener to the bus.
// Call once after InitBus has run.
func BusListeners() {
	b := x.EventBus()

	subscribe := func(event string, l bus.Listener) {
		if err := b.Subscribe(event, l); err != nil {
			errMsg := "bus: subscribe failed: " + err.Error()
			x.Logger().Error(errMsg)
			utils.PrintErr(errMsg)
		}
	}

	userRepository := x.AppMust[*users.UserRepository]()
	auditLogRepository := x.AppMust[*audit_logs.AuditLogRepository]()

	subscribe(events.EventUserLoggedIn, audit_logs_listeners.NewUserLoggedInListener(x.Notify(), auditLogRepository))
	subscribe(events.EventUserLoggedIn, notifications_listeners.NewUserLoggedInListener(x.Notify(), userRepository))
}

// taskHandlers register task handlers for tasks and
func taskHandlers() map[string]asynq.Handler {
	return map[string]asynq.Handler{
		// xnotify tasks
		xnotify.TaskType: asynqNotify.NewHandler(x.Notify()),
		// bus tasks
		x.EventBus().Prefix(): asynqBus.NewAsynqHandler(x.EventBus()),
	}
}
