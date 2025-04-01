package main

import (
	"fmt"
	"myApp/config"
	"myApp/model"
	"myApp/pkg/logger"
	"myApp/pkg/redis"
	"myApp/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化日志
	logger.Init(logger.Config{
		Level:      "info",
		FilePath:   "./logs/app.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
		Console:    true,
	})

	// 记录应用启动日志
	logger.Info("应用启动中",
		zap.String("mode", config.Conf.Server.Mode),
	)

	// 初始化数据库
	model.InitDB()

	// 初始化Redis
	redis.InitRedis()

	// 设置Gin运行模式
	if config.Conf.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin实例
	r := gin.New()
	r.Use(gin.Recovery())

	// 初始化路由
	router.SetupRouter(r)

	// 启动HTTP服务
	fmt.Printf("\n🚀 服务端启动成功，监听端口 %d\n", config.Conf.Server.Port)
	if err := r.Run(fmt.Sprintf(":%d", config.Conf.Server.Port)); err != nil {
		panic(fmt.Sprintf("服务启动失败: %v", err))
	}
}
