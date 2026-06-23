package load

import (
	"context"
	"errors"

	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/hooks"
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/pkg/inertia"
	"github.com/imohamedsheta/xapp/resources"
	"github.com/imohamedsheta/xioc"
)

func InitInertia(c *xioc.Container) {
	flash := initFlashProvider(c)
	initInertiaInstance(c, flash)
}

// ---------------- Flash Provider ----------------
func initFlashProvider(c *xioc.Container) *inertia.InmemFlashProvider {
	flash := inertia.NewInmemFlashProvider()

	_ = xioc.Singleton(c, func(c *xioc.Container) (*inertia.InmemFlashProvider, error) {
		flash.SetSessionIDFunc(func(ctx context.Context) (string, error) {
			sessionId, ok := ctx.Value(string(enums.ContextKeySessionId)).(string)
			if !ok || sessionId == "" {
				return "", nil
			}
			return sessionId, nil
		})
		return flash, nil
	})

	return flash
}

// ---------------- Inertia Instance ----------------
func initInertiaInstance(c *xioc.Container, flash *inertia.InmemFlashProvider) {
	err := xioc.Singleton(c, func(c *xioc.Container) (*inertia.Inertia, error) {
		// Pass the embedded FS and the path within it
		i := inertia.InitInertia(resources.ViewsFS, "views/root.html")
		if i == nil {
			return nil, errors.New("you should open hot reload for dev server using npm run dev")
		}
		hooks.InitInertiaHooks(i, flash)
		return i, nil
	})

	if err != nil {
		logError("Failed to load inertia module in the ioc container : " + err.Error())
	}
}

func logError(msg string) {
	utils.PrintErr(msg)
	x.Logger().Error(msg)
}
