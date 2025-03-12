package router

import (
	"myApp/handler"
	"myApp/middleware"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

// RegisterLandlordRoutes 注册房东相关路由
func RegisterLandlordRoutes(r *gin.RouterGroup) {
	// 创建仓库实例
	landlordRepo := repository.NewLandlordRepository()
	userRepo := repository.NewUserRepository()
	
	// 创建服务实例
	landlordService := service.NewLandlordService(landlordRepo, userRepo)
	
	// 创建处理器实例
	landlordHandler := handler.NewLandlordHandler(landlordService)
	
	// 房东路由组
	landlordRoutes := r.Group("/landlords")
	{
		// 需要认证的路由
		landlordRoutes.Use(middleware.JWTAuth())
		
		// 申请成为房东
		landlordRoutes.POST("", landlordHandler.CreateLandlord)
		
		// 获取房东个人资料
		landlordRoutes.GET("/profile", landlordHandler.GetLandlordProfile)
		
		// 更新房东信息
		landlordRoutes.PUT("/profile", landlordHandler.UpdateLandlord)
		
		// 管理员验证房东身份（需要管理员权限）
		landlordRoutes.PUT("/verify/:id", landlordHandler.VerifyLandlord)
	}
}