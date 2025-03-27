package middleware

import (
	"myApp/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimiter() gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(1*time.Second, 100, 100) // 每秒100个请求
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.Fail(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
