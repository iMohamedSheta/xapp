package x

import "github.com/imohamedsheta/xsocial"

func Socialite() *xsocial.Socialite {
	return AppMust[*xsocial.Socialite]()
}
