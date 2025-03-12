package router

import (
	"myApp/handler"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

// InitUserRouter 初始化用户相关路由
func InitUserRouter(r *gin.Engine) {
	// 创建用户数据仓库实例
	userRepo := repository.NewUserRepository()
	// 创建用户服务实例，注入数据仓库依赖
	userService := service.NewUserService(userRepo)
	// 创建用户处理器实例，注入服务依赖
	userHandler := handler.NewUserHandler(userService)
	
	// 创建用户路由组，所有用户相关接口都在/api/user路径下
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/register", userHandler.Register) // 用户注册接口
		userGroup.POST("/login", userHandler.Login)       // 用户登录接口
		userGroup.GET("/info", userHandler.GetUserInfo)   // 获取用户信息接口，需要JWT认证
	}
}
