package middleware

import (
	"myApp/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 结构化日志中间件
// 提供丰富的日志信息和结构化输出
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 初始化请求上下文日志
		ctxLogger := logger.WithContext(c)

		// 记录请求开始
		ctxLogger.Info("请求开始",
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// 记录请求时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(start)

		// 获取响应状态
		status := c.Writer.Status()

		// 根据状态码确定日志级别
		if status >= 500 {
			ctxLogger.Error("请求完成",
				zap.Int("status", status),
				zap.Duration("duration", duration),
				zap.Int("size", c.Writer.Size()),
				zap.String("errors", c.Errors.String()),
			)
		} else if status >= 400 {
			ctxLogger.Warn("请求完成",
				zap.Int("status", status),
				zap.Duration("duration", duration),
				zap.Int("size", c.Writer.Size()),
				zap.String("errors", c.Errors.String()),
			)
		} else {
			ctxLogger.Info("请求完成",
				zap.Int("status", status),
				zap.Duration("duration", duration),
				zap.Int("size", c.Writer.Size()),
			)
		}
	}
}
