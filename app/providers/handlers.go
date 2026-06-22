package providers

import (
	"github.com/imohamedsheta/xapp/app/http/handler"
	"github.com/imohamedsheta/xapp/app/modules/audit_logs"
	"github.com/imohamedsheta/xapp/app/modules/identity/auth"
	"github.com/imohamedsheta/xapp/app/modules/identity/tenants"
	"github.com/imohamedsheta/xapp/app/modules/identity/users"
	"github.com/imohamedsheta/xapp/app/modules/settings"
	"github.com/imohamedsheta/xioc"
)

// Main register function for all providers
func RegisterProviders(c *xioc.Container) {
	RegisterMiddlewares(c)
	RegisterRepositories(c)
	RegisterServices(c)
	RegisterActions(c)
	RegisterHandlers(c)
	RegisterPolicies(c)
}

// Register Services in Container
func RegisterServices(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, auth.NewJwtService))
	panicIfError(xioc.SingletonLazy(c, auth.NewPermissionService))
	panicIfError(xioc.SingletonLazy(c, users.NewUserService))
	panicIfError(xioc.SingletonLazy(c, tenants.NewTenantService))
	panicIfError(xioc.SingletonLazy(c, auth.NewAuthService))
}

// Register Actions in Container
func RegisterActions(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, auth.NewAuthAction))
	panicIfError(xioc.SingletonLazy(c, settings.NewSettingAction))
}

func RegisterPolicies(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, settings.NewSettingPolicy))
}

// Register Repositories in Container
func RegisterRepositories(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, settings.NewSettingRepository))
	panicIfError(xioc.SingletonLazy(c, auth.NewPermissionRepository))
	panicIfError(xioc.SingletonLazy(c, users.NewUserRepository))
	panicIfError(xioc.SingletonLazy(c, tenants.NewTenantRepository))
	panicIfError(xioc.SingletonLazy(c, audit_logs.NewAuditLogRepository))
}

// Register Handlers in Container
func RegisterHandlers(c *xioc.Container) {
	panicIfError(xioc.SingletonLazy(c, handler.NewBaseHandler))
	panicIfError(xioc.SingletonLazy(c, auth.NewAuthHandler))
	panicIfError(xioc.SingletonLazy(c, settings.NewSettingsHandler))
	panicIfError(xioc.SingletonLazy(c, users.NewAccountHandler))
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
