package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/pkg/inertia"
)

func FlashToast(c *gin.Context, toast ...map[string]any) error {
	return x.Flash().Flash(c, inertia.Props{
		"toast": toast,
	})
}

func FlashBack(c *gin.Context, i *inertia.Inertia, status int, toast ...map[string]any) error {
	err := x.Flash().Flash(c, inertia.Props{
		"toast": toast,
	})

	if err != nil {
		return err
	}

	i.Back(c, 303)
	return nil
}
