package main

import (
	"fmt"
	"myApp/config"
	"myApp/model"
	"myApp/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库
	model.InitDB()

	// 创建Gin实例
	r := gin.Default()
	router.SetupRouter(r)

	// 启动HTTP服务
	fmt.Printf("\n🚀 服务端启动成功，监听端口 %d\n", config.Conf.Server.Port)
	if err := r.Run(fmt.Sprintf(":%d", config.Conf.Server.Port)); err != nil {
		panic(fmt.Sprintf("服务启动失败: %v", err))
	}
}