package router

import (
	"myApp/handler"
	"myApp/middleware"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

// InitViewingRouter 初始化预约看房相关路由
func InitViewingRouter(r *gin.Engine) {
	// 创建预约看房数据仓库实例
	viewingRepo := repository.NewViewingRepository()
	// 创建预约看房服务实例，注入数据仓库依赖
	viewingService := service.NewViewingService(viewingRepo)

	// 创建房源数据仓库实例
	houseRepo := repository.NewHouseRepository()
	// 创建房源服务实例，注入数据仓库依赖
	houseService := service.NewHouseService(houseRepo)

	// 创建预约看房处理器实例，注入服务依赖
	viewingHandler := handler.NewViewingHandler(viewingService, houseService)

	// 创建预约看房路由组，所有预约看房相关接口都在/api/viewing路径下
	viewingGroup := r.Group("/api/viewing")
	// 所有预约看房接口都需要认证，添加JWT中间件
	viewingGroup.Use(middleware.JWTAuth())
	{
		viewingGroup.POST("/create", viewingHandler.CreateViewing)            // 创建预约看房
		viewingGroup.GET("/:id", viewingHandler.GetViewing)                   // 获取预约看房详情
		viewingGroup.GET("/user", viewingHandler.GetUserViewings)             // 获取用户的所有预约看房
		viewingGroup.GET("/house/:house_id", viewingHandler.GetHouseViewings) // 获取房源的所有预约看房
		viewingGroup.PUT("/confirm/:id", viewingHandler.ConfirmViewing)       // 确认预约看房
		viewingGroup.PUT("/complete/:id", viewingHandler.CompleteViewing)     // 完成预约看房
		viewingGroup.PUT("/cancel/:id", viewingHandler.CancelViewing)         // 取消预约看房
	}
}
