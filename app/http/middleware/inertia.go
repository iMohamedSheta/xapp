package middleware

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/app/x"
)

// InertiaMiddlewareWithErrorHandler wraps a HandlerFuncWithError with Inertia’s middleware.
// It also records the original handler name into gin.Context for debugging.
func InertiaMiddlewareWithErrorHandler(next HandlerFuncWithError) gin.HandlerFunc {
	i := x.Inertia()

	// extract the real function name once
	fnName := runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name()
	// optional: shorten it
	short := shortActionName(fnName)

	return func(c *gin.Context) {
		// save for later (debug tab)
		c.Set("route_action_full", fnName)
		c.Set("route_action_short", short)

		var handlerError error

		// Build http.Handler chain using inertia.Middleware
		h := i.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Execute the actual handler
			handlerError = next(c)

			// If there's an error, handle it with the global error handler
			if handlerError != nil {
				ErrorHandler(c, handlerError, true)
			}
		}))

		h.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

// same helper you already had
func shortActionName(full string) string {
	parts := strings.Split(full, ".")
	return parts[len(parts)-1]
}
