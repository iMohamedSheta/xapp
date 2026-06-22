package bootstrap

import (
	"github.com/go-playground/validator/v10"
	"github.com/imohamedsheta/xapp/bootstrap/load"
	"github.com/imohamedsheta/xapp/bootstrap/support"
	"github.com/imohamedsheta/xioc"
	"github.com/imohamedsheta/xnotify"

	"github.com/joho/godotenv"
)

type AppBuilder struct {
	envFile   string
	actions   []func()
	container *xioc.Container
}

// NewAppBuilder starts a new application builder.
func NewAppBuilder(envFile string) *AppBuilder {
	container := xioc.New()
	xioc.SetAppContainer(container)
	return &AppBuilder{
		envFile:   envFile,
		container: container,
	}
}

// Then inserts a custom step into the builder chain.
func (b *AppBuilder) Then(fn func(b *AppBuilder)) *AppBuilder {
	fn(b)
	return b
}

// LoadConfig loads and registers the xfig config singleton.
func (b *AppBuilder) LoadConfig() *AppBuilder {
	load.InitConfig(b.container)
	return b
}

func (b *AppBuilder) MustLoadEnvFile() *AppBuilder {
	b.loadEnvFile(true)
	return b
}

func (b *AppBuilder) LoadEventBus() *AppBuilder {
	load.InitEventBus(b.container)
	return b
}

func (b *AppBuilder) LoadEnvFile() *AppBuilder {
	b.loadEnvFile(false)
	return b
}

func (b *AppBuilder) LoadLogger() *AppBuilder {
	load.InitLogger(b.container)
	return b
}

func (b *AppBuilder) LoadStorage() *AppBuilder {
	load.InitStorage(b.container)
	return b
}

func (b *AppBuilder) LoadInertia() *AppBuilder {
	load.InitInertia(b.container)
	return b
}

func (b *AppBuilder) LoadDatabase() *AppBuilder {
	load.InitDatabase(b.container)
	return b
}

// LoadValidator registers custom validation rules into the validator singleton.
func (b *AppBuilder) LoadValidator(rules map[string]validator.FuncCtx) *AppBuilder {
	load.InitValidator(b.container, rules)
	return b
}

func (b *AppBuilder) LoadRedisCache() *AppBuilder {
	load.InitRedisCache(b.container)
	return b
}

func (b *AppBuilder) LoadSocialite() *AppBuilder {
	load.InitSocialite(b.container)
	return b
}

func (b *AppBuilder) LoadRedisQueue() *AppBuilder {
	load.InitRedisQueue(b.container)
	return b
}

// LoadNotify registers notification channel handlers into the notify singleton.
func (b *AppBuilder) LoadNotify(channels map[string]xnotify.ChannelHandler) *AppBuilder {
	load.InitNotify(b.container, channels)
	return b
}

func (b *AppBuilder) LoadWebsocketServer() *AppBuilder {
	load.InitWebsocketServer(b.container)
	return b
}

func (b *AppBuilder) LoadXErr() *AppBuilder {
	load.InitXErr(b.container)
	return b
}

func (b *AppBuilder) loadEnvFile(must bool) {
	if err := godotenv.Load(b.envFile); err != nil {
		if must {
			support.PrintErr("Error loading environment file: " + err.Error())
			panic("can't load environment file: " + b.envFile + ": " + err.Error())
		}
	}
	support.PrintSuccess("Loaded environment file: " + b.envFile)
}

// Boot finalizes the IoC container. Pass your RegisterServiceProviders func here.
func (b *AppBuilder) Boot(registerProviders func(*xioc.Container)) *AppBuilder {
	registerProviders(b.container)
	b.container.Bootstrap()
	return b
}
