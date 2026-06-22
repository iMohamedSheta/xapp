package x

import "github.com/iMohamedSheta/xerr"

func XErr() *xerr.ErrorHandler {
	return AppMust[*xerr.ErrorHandler]()
}
