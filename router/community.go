package router

import (
	"myApp/handler"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

func InitCommunityRouter(r *gin.Engine) {
	communityRepo := repository.NewCommunityRepository()
	communityService := service.NewCommunityService(communityRepo)
	communityHandler := handler.NewCommunityHandler(communityService)
	communityGroup := r.Group("/api/community")
	{
		communityGroup.POST("/create", communityHandler.CreateCommunity)
		communityGroup.GET("/list", communityHandler.GetCommunityList)
		communityGroup.GET("/:id", communityHandler.GetCommunity)
	}
}
