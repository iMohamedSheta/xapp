package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/app/domain/enums"
)

func GetRequestId(c *gin.Context) string {
	return c.GetString(enums.ContextKeyRequestId.String())
}
