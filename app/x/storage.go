package x

import "github.com/imohamedsheta/xdisk"

func Storage() *xdisk.Storage {
	return AppMust[*xdisk.Storage]()
}
