package router

import (
	"myApp/handler"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

// InitFavoriteRouter 初始化收藏相关路由
func InitFavoriteRouter(r *gin.Engine) {
	// 创建收藏数据仓库实例
	favoriteRepo := repository.NewFavoriteRepository()
	// 创建收藏服务实例，注入数据仓库依赖
	favoriteService := service.NewFavoriteService(favoriteRepo)
	// 创建收藏处理器实例，注入服务依赖
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	// 创建收藏路由组，所有收藏相关接口都在/api/favorite路径下
	favoriteGroup := r.Group("/api/favorite")
	{
		// 所有收藏接口都需要认证
		favoriteGroup.POST("/add", favoriteHandler.AddFavorite)                // 添加收藏
		favoriteGroup.DELETE("/:id", favoriteHandler.RemoveFavorite)          // 删除收藏
		favoriteGroup.GET("/list", favoriteHandler.GetUserFavorites)         // 获取用户的所有收藏
		favoriteGroup.POST("/toggle/:house_id", favoriteHandler.ToggleFavorite) // 切换收藏状态
		favoriteGroup.GET("/check/:house_id", favoriteHandler.CheckFavorite)   // 检查是否已收藏
	}
}