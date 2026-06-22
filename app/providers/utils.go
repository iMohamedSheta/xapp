package providers

import (
	"github.com/imohamedsheta/xapp/app/x"
)

// Helper that create error logs error message for the given module name and error
func logBindErr(module string, err error) {
	if err != nil {
		x.Logger().Error("Failed to load " + module + " module in the ioc container : " + err.Error())
	}
}
