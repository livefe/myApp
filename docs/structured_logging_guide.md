# 结构化日志使用指南

## 概述

本项目已集成结构化日志系统，基于 zap 日志库实现。结构化日志相比传统日志有以下优势：

1. **结构化数据**：以 JSON 格式输出，便于日志分析和搜索
2. **上下文信息**：自动包含请求 ID、用户 ID 等上下文信息
3. **性能优异**：zap 是目前 Go 生态中性能最好的日志库之一
4. **级别控制**：支持不同级别的日志记录（Debug, Info, Warn, Error, Fatal）
5. **字段扩展**：支持添加自定义字段，提供更丰富的日志信息

## 日志级别使用指南

- **Debug**：详细的调试信息，仅在开发环境使用
- **Info**：常规操作信息，记录系统正常运行状态
- **Warn**：潜在问题或异常情况，但不影响系统正常运行
- **Error**：错误信息，影响功能但不影响系统运行
- **Fatal**：严重错误，导致系统无法继续运行

## 在服务层使用结构化日志

```go
// 导入日志包
import (
    "myApp/pkg/logger"
    "go.uber.org/zap"
)

// 在服务方法中使用
func (s *someService) SomeOperation(param string) error {
    // 记录操作开始
    logger.Info("操作开始", 
        zap.String("param", param),
        zap.String("operation", "SomeOperation"),
    )
    
    // 执行业务逻辑...
    result, err := s.repo.DoSomething(param)
    
    if err != nil {
        // 记录错误
        logger.Error("操作失败", 
            zap.String("param", param),
            zap.String("operation", "SomeOperation"),
            zap.Error(err),
        )
        return err
    }
    
    // 记录操作成功
    logger.Info("操作成功", 
        zap.String("param", param),
        zap.String("operation", "SomeOperation"),
        zap.Any("result", result),
    )
    
    return nil
}
```

## 在处理函数中使用结构化日志

```go
// 导入日志包
import (
    "myApp/pkg/logger"
    "go.uber.org/zap"
)

// 在处理函数中使用上下文日志
func (h *SomeHandler) HandleRequest(c *gin.Context) {
    // 获取请求参数
    id := c.Param("id")
    
    // 记录请求信息
    logger.ContextInfo(c, "处理请求开始", 
        zap.String("id", id),
        zap.String("handler", "HandleRequest"),
    )
    
    // 执行业务逻辑
    result, err := h.service.SomeOperation(id)
    
    if err != nil {
        // 记录错误
        logger.ContextError(c, "处理请求失败", 
            zap.String("id", id),
            zap.String("handler", "HandleRequest"),
            zap.Error(err),
        )
        // 返回错误响应...
        return
    }
    
    // 记录成功信息
    logger.ContextInfo(c, "处理请求成功", 
        zap.String("id", id),
        zap.String("handler", "HandleRequest"),
    )
    
    // 返回成功响应...
}
```

## 最佳实践

1. **使用正确的日志级别**：根据信息的重要性选择合适的日志级别
2. **添加上下文信息**：在日志中包含足够的上下文信息，如操作类型、关键参数、用户ID等
3. **记录关键节点**：在操作的开始、结束以及关键步骤记录日志
4. **错误日志详细化**：记录错误时，包含详细的错误信息和上下文
5. **避免敏感信息**：不要在日志中记录密码、令牌等敏感信息
6. **使用结构化字段**：使用zap.String(), zap.Int()等方法添加结构化字段，而不是使用格式化字符串
7. **在处理函数中使用上下文日志**：使用logger.ContextInfo(), logger.ContextError()等方法，自动包含请求ID等信息

## 日志配置

日志配置在应用启动时在`main.go`中设置：

```go
import (
    "myApp/pkg/logger"
    "myApp/config"
)

// 在应用启动时初始化
// 注意：需要先调用config.InitConfig()初始化配置
logger.Init() // 直接从config.Conf获取日志配置
```

可以根据环境需要调整这些配置参数。