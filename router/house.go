package router

import (
	"myApp/handler"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

// InitHouseRouter 初始化房源相关路由
func InitHouseRouter(r *gin.Engine) {
	// 创建房源数据仓库实例
	houseRepo := repository.NewHouseRepository()
	// 创建房源服务实例，注入数据仓库依赖
	houseService := service.NewHouseService(houseRepo)
	// 创建房源处理器实例，注入服务依赖
	houseHandler := handler.NewHouseHandler(houseService)

	// 创建房源路由组，所有房源相关接口都在/api/house路径下
	houseGroup := r.Group("/api/house")
	{
		// 公开接口，不需要认证
		houseGroup.GET("/list", houseHandler.GetAllHouses) // 获取房源列表
		houseGroup.GET("/:id", houseHandler.GetHouse)      // 获取房源详情

		// 需要认证的接口，需要JWT认证
		authorizedGroup := houseGroup.Group("/")
		{
			authorizedGroup.POST("/create", houseHandler.CreateHouse)         // 创建房源
			authorizedGroup.PUT("/:id", houseHandler.UpdateHouse)             // 更新房源
			authorizedGroup.DELETE("/:id", houseHandler.DeleteHouse)          // 删除房源
			authorizedGroup.GET("/landlord", houseHandler.GetLandlordHouses) // 获取房东的所有房源
		}
	}
}