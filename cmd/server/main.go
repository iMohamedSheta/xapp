package main

import (
	registers "github.com/imohamedsheta/xapp/app"
	"github.com/imohamedsheta/xapp/app/domain/tasks/scheduler"
	apphttp "github.com/imohamedsheta/xapp/app/http"
	"github.com/imohamedsheta/xapp/bootstrap"
	"github.com/imohamedsheta/xapp/bootstrap/support"
)

func main() {
	bootstrap.NewAppBuilder(".env").
		MustLoadEnvFile().
		LoadConfig().
		LoadLogger().
		LoadDatabase().
		LoadStorage().
		LoadValidator(registers.ValidationRules()).
		LoadRedisCache().
		LoadRedisQueue().
		LoadWebsocketServer().
		LoadNotify(registers.NotifyChannels()).
		LoadInertia().
		LoadSocialite().
		LoadEventBus().
		LoadXErr().
		Boot(registers.ServiceProviders)

	// Register WebSocket channel policies now that WS server is booted.
	registers.WebSocketChannels()

	support.SafeGo(func() { bootstrap.RunWorker(registers.TaskHandlers(), false) })
	support.SafeGo(func() { bootstrap.RunWebsocket(registers.TaskHandlers(), registers.WebsocketHooks()) })
	support.SafeGo(func() { bootstrap.RunScheduler(scheduler.RegisterSchedule) })

	bootstrap.RunHttp(apphttp.RegisterRoutes())
}
