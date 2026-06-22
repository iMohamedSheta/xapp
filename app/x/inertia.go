package x

import "github.com/imohamedsheta/xapp/pkg/inertia"

func Inertia() *inertia.Inertia {
	return AppMust[*inertia.Inertia]()
}

func Flash() *inertia.InmemFlashProvider {
	return AppMust[*inertia.InmemFlashProvider]()
}
