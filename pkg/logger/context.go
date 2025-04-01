package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ContextKey 用于在gin.Context中存储日志相关的键
type ContextKey string

const (
	// RequestIDKey 请求ID的键
	RequestIDKey ContextKey = "request_id"
	// LoggerKey 日志实例的键
	LoggerKey ContextKey = "logger"
)

// GetRequestID 从上下文中获取请求ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(string(RequestIDKey)); exists {
		return requestID.(string)
	}
	return ""
}

// GetContextLogger 从上下文中获取日志实例
func GetContextLogger(c *gin.Context) *zap.Logger {
	if logger, exists := c.Get(string(LoggerKey)); exists {
		return logger.(*zap.Logger)
	}
	return Logger
}

// SetContextLogger 设置上下文日志实例
func SetContextLogger(c *gin.Context, logger *zap.Logger) {
	c.Set(string(LoggerKey), logger)
}

// WithContext 创建带有上下文信息的日志实例
func WithContext(c *gin.Context) *zap.Logger {
	// 获取或生成请求ID
	requestID := GetRequestID(c)
	if requestID == "" {
		requestID = uuid.New().String()
		c.Set(string(RequestIDKey), requestID)
	}

	// 创建带有请求信息的日志实例
	logger := Logger.With(
		zap.String("request_id", requestID),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("client_ip", c.ClientIP()),
	)

	// 如果有用户ID，添加到日志中
	if userID, exists := c.Get("user_id"); exists {
		logger = logger.With(zap.Any("user_id", userID))
	}

	// 保存到上下文中
	SetContextLogger(c, logger)

	return logger
}

// ContextInfo 输出Info级别的上下文日志
func ContextInfo(c *gin.Context, msg string, fields ...zap.Field) {
	GetContextLogger(c).Info(msg, fields...)
}

// ContextError 输出Error级别的上下文日志
func ContextError(c *gin.Context, msg string, fields ...zap.Field) {
	GetContextLogger(c).Error(msg, fields...)
}

// ContextWarn 输出Warn级别的上下文日志
func ContextWarn(c *gin.Context, msg string, fields ...zap.Field) {
	GetContextLogger(c).Warn(msg, fields...)
}

// ContextDebug 输出Debug级别的上下文日志
func ContextDebug(c *gin.Context, msg string, fields ...zap.Field) {
	GetContextLogger(c).Debug(msg, fields...)
}
