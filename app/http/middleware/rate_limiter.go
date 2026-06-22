package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/app/x"
	limiter "github.com/ulule/limiter/v3"
	limiterGinMiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

// in memory rate limiter maybe will use redis store in future
func RateLimiter() gin.HandlerFunc {
	globalRateLimiter := x.Config().GetString("app.global_rate_limit", "10000-M") // Default to 100 requests per minute
	rate, _ := limiter.NewRateFromFormatted(globalRateLimiter)

	store := memory.NewStore()

	instance := limiter.New(store, rate)

	// Create a new Gin rate limiter middleware using the limiter instance
	middleware := limiterGinMiddleware.NewMiddleware(instance)

	return middleware
}
