# 结构化日志包

这个包提供了结构化日志功能，基于zap日志库实现，支持以下特性：

- 多级别日志（Debug, Info, Warn, Error, Fatal）
- JSON格式的结构化日志输出
- 日志文件自动轮转（基于大小、时间）
- 请求上下文跟踪（请求ID、用户ID等）
- 支持同时输出到文件和控制台
- 支持添加自定义字段

## 使用方法

### 1. 初始化日志

在应用启动时初始化日志：

```go
import (
    "myApp/pkg/logger"
    "myApp/config"
)

// 在应用启动时初始化
// 注意：需要先调用config.InitConfig()初始化配置
logger.InitLogger() // 直接从config.Conf获取日志配置
```

### 2. 基本日志记录

```go
// 记录不同级别的日志
logger.Debug("这是一条调试日志")
logger.Info("这是一条信息日志")
logger.Warn("这是一条警告日志")
logger.Error("这是一条错误日志")

// 带有额外字段的日志
logger.Info("用户登录", zap.String("username", "张三"), zap.Int("user_id", 123))

// 格式化日志
logger.Infof("用户 %s (ID: %d) 登录成功", "张三", 123)

// 带有错误信息的日志
err := someFunction()
if err != nil {
    logger.WithError(err).Error("操作失败")
}

// 带有多个字段的日志
fields := map[string]interface{}{
    "user_id": 123,
    "action": "login",
    "ip": "192.168.1.1",
}
logger.WithFields(fields).Info("用户活动")
```

### 3. 请求上下文日志

在中间件中初始化请求上下文日志：

```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 初始化请求上下文日志
        logger.WithContext(c)
        
        start := time.Now()
        c.Next()
        duration := time.Since(start)
        
        // 记录请求处理信息
        logger.ContextInfo(c, "请求完成",
            zap.Int("status", c.Writer.Status()),
            zap.Duration("duration", duration),
        )
    }
}
```

在处理函数中使用上下文日志：

```go
func SomeHandler(c *gin.Context) {
    // 使用上下文日志记录信息
    logger.ContextInfo(c, "处理请求", zap.String("param", c.Param("id")))
    
    // 处理业务逻辑...
    
    if err != nil {
        logger.ContextError(c, "处理失败", zap.Error(err))
        // 返回错误响应...
        return
    }
    
    logger.ContextInfo(c, "处理成功")
    // 返回成功响应...
}
```