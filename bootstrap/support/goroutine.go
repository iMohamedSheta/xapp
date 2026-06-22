package support

import (
	"fmt"
	"runtime/debug"

	"github.com/imohamedsheta/xapp/app/x"
)

// SafeGo runs fn in a goroutine, recovering from panics and logging them.
func SafeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				x.Logger().Error(fmt.Sprintf(
					"Panic in goroutine: %v\nStack: %s",
					r,
					debug.Stack(),
				))
			}
		}()
		fn()
	}()
}
