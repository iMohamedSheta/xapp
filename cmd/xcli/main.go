package main

import (
	"fmt"
	"os"
	"slices"

	registers "github.com/imohamedsheta/xapp/app"
	"github.com/imohamedsheta/xapp/bootstrap"
	"github.com/imohamedsheta/xcli"
)

var cliOnlyCommands = []string{
	"inspire",
	"seed",
}

func main() {
	// 1. Determine environment file (load .env.xcli if it exists, fallback to .env)
	envFile := ".env"
	if _, err := os.Stat(".env.xcli"); err == nil {
		envFile = ".env.xcli"
	}

	// 2. Peek at the requested command argument
	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	// check if command is cli only
	if xcli.IsCliOnly(cmd) || slices.Contains(cliOnlyCommands, cmd) {
		// Minimal boot: Config + Logger only (no database, no redis connections needed)
		bootstrap.NewAppBuilder(envFile).
			LoadEnvFile().
			LoadConfig().
			LoadLogger().
			Boot(registers.ServiceProviders)
	} else {
		// Full boot: Load database, redis, websocket, notifications, etc.
		bootstrap.NewAppBuilder(envFile).
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
			LoadXErr().
			Boot(registers.ServiceProviders)
	}

	// 3. Instantiate xcli
	x := xcli.New()

	// 4. Register custom application commands
	// You can load them from registers, or declare them directly:
	x.Register(registers.Commands()...)

	// 5. Run the CLI
	if err := x.Execute(); err != nil {
		fmt.Printf("\033[31m❌ Command failed: %v\033[0m\n", err)
		os.Exit(1)
	}
}
