package middleware

import (
	"fmt"

	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xapp/app/x"
)

func RecoveryWithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				x.Logger().Error(fmt.Sprintf(
					"\n🚨 Panic Recovered 🚨\nMethod: %s\nEndpoint: %s\nError: %v\n\nStack Trace:\n%s\n\n",
					c.Request.Method,
					c.Request.URL.Path,
					r,
					debug.Stack(),
				))

				if utils.IsDebug() && !strings.HasPrefix(c.FullPath(), "/api/") {
					x.XErr().HandleError(c.Writer, c.Request, r)
				}

				ErrorHandler(c, xerr.New("Panic error unexpected error just happend", enums.XErrServerError, fmt.Errorf("%v", r)), true)
			}
		}()
		c.Next()
	}
}
