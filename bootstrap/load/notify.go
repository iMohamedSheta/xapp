package load

import (
	"os"

	"github.com/hibiken/asynq"
	"github.com/imohamedsheta/xapp/app/domain/adapters"
	"github.com/imohamedsheta/xapp/app/domain/tasks"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xioc"
	"github.com/imohamedsheta/xnotify"
	asynqNotify "github.com/imohamedsheta/xnotify/asynq"
)

func InitNotify(c *xioc.Container, channelsHandlers map[string]xnotify.ChannelHandler) {
	err := xioc.Singleton(c, func(c *xioc.Container) (*xnotify.Notify, error) {
		zapLogger := x.Logger()
		logger := adapters.NewLoggerAdapter(zapLogger)
		queue, err := xioc.AppMake[*x.QueueClient]()
		var backend xnotify.QueueBackend
		if err != nil {
			x.Logger().Error("Failed to load queue client in the ioc container: " + err.Error())
			backend = nil
		} else {
			backend = asynqNotify.New(queue.Client, asynq.Queue(tasks.QueueNotifications), asynq.MaxRetry(3))
		}

		notify := xnotify.New(logger, backend)
		notify.RegisterChannels(channelsHandlers)
		return notify, nil
	})

	if err != nil {
		errMsg := "Failed to load notify module as singleton in the ioc container: " + err.Error()
		x.Logger().Error(errMsg)
		utils.PrintErr(errMsg)
		os.Exit(1)
	}
}
