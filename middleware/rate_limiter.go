package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

func RateLimiter() gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(1*time.Second, 100, 100) // 每秒100个请求
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, 
				gin.H{"error": "请求过于频繁，请稍后再试"})
			return
		}
		c.Next()
	}
}