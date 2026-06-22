package x

import (
	"fmt"

	"github.com/imohamedsheta/xapp/pkg/bus"
	"github.com/imohamedsheta/xnotify"
)

func Notify() *xnotify.Notify {
	n, err := app[*xnotify.Notify]()
	if err != nil {
		Logger().Error(fmt.Sprintf("Notify can't be resolved: %s", err.Error()))
		return nil
	}
	return n
}

func EventBus() *bus.Bus {
	b, err := app[*bus.Bus]()
	if err != nil {
		Logger().Error(fmt.Sprintf("EventBus can't be resolved: %s", err.Error()))
		return nil
	}
	return b
}
