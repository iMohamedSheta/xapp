package utils

import "github.com/gin-gonic/gin"

// Problem with this it doesn't add middlewares to route as its
func InternalRedirect(to string, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = to
		router.HandleContext(c)
	}
}

func Redirect(to string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(301, to)
		c.Abort()
	}
}
