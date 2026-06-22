package load

import (
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xioc"
)

func InitXErr(c *xioc.Container) {
	err := xioc.Singleton(c, func(c *xioc.Container) (*xerr.ErrorHandler, error) {
		return xerr.NewErrorHandler(&xerr.Config{
			ShowSourceCode: true,
			MaxFrames:      50,
			Environment:    x.Config().GetString("app.env", "dev"),
			DebugMode:      x.Config().GetBool("app.debug", true),
			SkipFrames:     0,
			SkipLibrary:    true,
		}), nil
	})

	if err != nil {
		x.Logger().Error("Failed to load xerr module in the ioc container : " + err.Error())
	}
}
