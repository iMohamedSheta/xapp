package load

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xioc"
	"github.com/imohamedsheta/xvalid"
	fileupload "github.com/imohamedsheta/xvalid/fileupload"
)

func InitValidator(c *xioc.Container, registeredRules map[string]validator.FuncCtx) {
	err := xioc.Singleton(c, func(c *xioc.Container) (*xvalid.Validator, error) {
		v := validator.New()

		for tag, rule := range registeredRules {
			if err := v.RegisterValidationCtx(tag, rule); err != nil {
				return nil, err
			}
		}

		fileupload.RegisterFileHeaderType(v)
		fileupload.RegisterFileValidations(v)

		return xvalid.New(v), nil
	})

	if err != nil {
		errMsg := "Failed to load validator module in the ioc container : " + err.Error()
		x.Logger().Error(errMsg)
		utils.PrintErr(errMsg)
		os.Exit(1)
	}
}
