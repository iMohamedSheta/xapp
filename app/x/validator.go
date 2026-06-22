package x

import (
	"github.com/imohamedsheta/xvalid"
)

func Validator() *xvalid.Validator {
	v, err := app[*xvalid.Validator]()
	if err != nil {
		panic(err)
	}
	return v
}
