package x

import (
	"fmt"

	"github.com/imohamedsheta/xws"
)

func WS() *xws.Server {
	ws, err := app[*xws.Server]()
	if err != nil {
		Logger().Error(fmt.Sprintf("WS Dependency Container Error: %s", err.Error()))
		return nil
	}

	return ws
}
