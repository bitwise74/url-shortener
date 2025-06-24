// Package middleware contains all custom middleware used by the server
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

func RateLimiter() gin.HandlerFunc {
	var limiter *rate.Limiter

	switch viper.GetString("rate_limiter_mode") {
	case "strict":
		limiter = rate.NewLimiter(1, 10)
	case "moderate":
		limiter = rate.NewLimiter(2, 20)
	case "relaxed":
		limiter = rate.NewLimiter(5, 100)
	case "disabled":
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Rate limit exceeded",
			})
			c.Abort()
		}
	}
}
