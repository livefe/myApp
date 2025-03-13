package router

import (
	"myApp/handler"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

// InitLandlordRouter 初始化房东相关路由
func InitLandlordRouter(r *gin.Engine) {
	// 创建房东数据仓库实例
	landlordRepo := repository.NewLandlordRepository()
	// 创建用户数据仓库实例
	userRepo := repository.NewUserRepository()

	// 创建房东服务实例，注入数据仓库依赖
	landlordService := service.NewLandlordService(landlordRepo, userRepo)

	// 创建房东处理器实例，注入服务依赖
	landlordHandler := handler.NewLandlordHandler(landlordService)

	// 创建房东路由组，所有房东相关接口都在/api/landlord路径下
	landlordGroup := r.Group("/api/landlord")
	{
		// 所有房东接口都需要认证
		landlordGroup.POST("/create", landlordHandler.CreateLandlord)         // 申请成为房东
		landlordGroup.GET("/profile", landlordHandler.GetLandlordProfile)    // 获取房东个人资料
		landlordGroup.PUT("/profile", landlordHandler.UpdateLandlord)        // 更新房东信息
		landlordGroup.PUT("/verify/:id", landlordHandler.VerifyLandlord)     // 管理员验证房东身份（需要管理员权限）
	}
}
