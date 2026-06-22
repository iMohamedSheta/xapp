package x

import (
	"fmt"

	"github.com/imohamedsheta/xioc"
)

// App resolves a required dependency. Panics if missing — use for
// dependencies the app genuinely cannot function without.
func AppMust[T any]() T {
	return xioc.AppMust[T]()
}

// App resolves an optional dependency. Returns (zero, false) if
// missing instead of panicking — use when the caller has a sane
// fallback and missing-ness is an expected, handleable case.
func App[T any]() (T, bool) {
	service, err := xioc.AppMake[T]()
	if err != nil {
		Logger().Error(fmt.Sprintf("AppMake[T] error: %v", err))
		return service, false
	}
	return service, true
}

func app[T any]() (service T, err error) {
	return xioc.AppMake[T]()
}
