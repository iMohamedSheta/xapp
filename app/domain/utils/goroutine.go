package utils

import (
	"fmt"
	"runtime/debug"

	"github.com/imohamedsheta/xapp/app/x"
)

// SafeGo is a wrapper for goroutines that will recover from panics and log them.
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
