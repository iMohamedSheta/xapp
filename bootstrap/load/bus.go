package load

import (
	"os"

	"github.com/hibiken/asynq"
	"github.com/imohamedsheta/xapp/app/domain/tasks"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/pkg/bus"
	bus_asynq_backend "github.com/imohamedsheta/xapp/pkg/bus/asynq_backned"
	"github.com/imohamedsheta/xioc"
)

func InitEventBus(c *xioc.Container) {
	err := xioc.Singleton(c, func(c *xioc.Container) (*bus.Bus, error) {
		queue, err := xioc.AppMake[*x.QueueClient]()
		var backend bus.QueueBackend
		if err != nil {
			x.Logger().Error("Failed to load queue client for bus: " + err.Error())
			backend = nil // Bus still works for sync listeners; queued ones error per-publish
		} else {
			backend = bus_asynq_backend.NewAsynqBackend(queue.Client, asynq.Queue(tasks.QueueEvents), asynq.MaxRetry(3))
		}

		return bus.New(backend, "bus:"), nil
	})

	if err != nil {
		errMsg := "Failed to load bus module as singleton in the ioc container: " + err.Error()
		x.Logger().Error(errMsg)
		utils.PrintErr(errMsg)
		os.Exit(1)
	}
}
