package providers

import (
	"github.com/imohamedsheta/xapp/app/http/handler"
	"github.com/imohamedsheta/xapp/app/modules/audit_logs"
	"github.com/imohamedsheta/xapp/app/modules/identity/auth"
	"github.com/imohamedsheta/xapp/app/modules/identity/tenants"
	"github.com/imohamedsheta/xapp/app/modules/identity/users"
	"github.com/imohamedsheta/xapp/app/modules/notifications"
	notifications_listeners "github.com/imohamedsheta/xapp/app/modules/notifications/listeners"
	"github.com/imohamedsheta/xapp/app/modules/settings"
	"github.com/imohamedsheta/xioc"
)

// Main register function for all providers
func RegisterProviders(c *xioc.Container) {
	RegisterMiddlewares(c)
	RegisterShared(c)
	//
	RegisterSettingsModule(c)
	RegisterAuditLogsModule(c)
	RegisterUsersModule(c)
	RegisterAuthModule(c)
	RegisterTenantsModule(c)
}

func RegisterModules(c *xioc.Container) {
	// identity bounded context
	RegisterAuthModule(c)
	RegisterUsersModule(c)
	RegisterTenantsModule(c)

	// notifications bounded context
	RegisterNotificationsModule(c)

	// audit logs bounded context
	RegisterAuditLogsModule(c)

	// settings bounded context
	RegisterSettingsModule(c)
}

func RegisterNotificationsModule(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, notifications.NewNotificationRepository))
	panicIfError(xioc.SingletonLazy(c, notifications_listeners.NewUserLoggedInListener))
}

func RegisterSettingsModule(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, settings.NewSettingRepository))
	panicIfError(xioc.SingletonLazy(c, settings.NewSettingAction))
	panicIfError(xioc.SingletonLazy(c, settings.NewSettingPolicy))
	panicIfError(xioc.SingletonLazy(c, settings.NewSettingsHandler))
}

func RegisterAuditLogsModule(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, audit_logs.NewAuditLogRepository))
}

func RegisterUsersModule(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, users.NewUserRepository))
	panicIfError(xioc.SingletonLazy(c, users.NewUserService))
	panicIfError(xioc.SingletonLazy(c, users.NewAccountHandler))
}

func RegisterAuthModule(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, auth.NewPermissionRepository))
	panicIfError(xioc.SingletonLazy(c, auth.NewPermissionService))
	panicIfError(xioc.SingletonLazy(c, auth.NewAuthService))
	panicIfError(xioc.SingletonLazy(c, auth.NewAuthAction))
	panicIfError(xioc.SingletonLazy(c, auth.NewAuthHandler))
	panicIfError(xioc.SingletonLazy(c, auth.NewJwtService))
}

func RegisterTenantsModule(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, tenants.NewTenantRepository))
	panicIfError(xioc.SingletonLazy(c, tenants.NewTenantService))
}

// Register Handlers in Container
func RegisterShared(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, handler.NewBaseHandler))
}

// Register Middlewares in Container
func RegisterMiddlewares(c *xioc.Container) {
	xioc.SingletonLazy(c, auth.NewAuthMiddleware)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
