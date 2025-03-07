package router

import (
	"myApp/handler/community"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

func InitCommunityRouter(r *gin.Engine) {
	communityRepo := repository.NewCommunityRepository()
	communityService := service.NewCommunityService(communityRepo)
	handler := community.NewCommunityHandler(communityService)
	communityGroup := r.Group("/api/community")
	{
		communityGroup.POST("/create", handler.CreateCommunity)
		communityGroup.GET("/list", handler.GetCommunityList)
		communityGroup.GET("/:id", handler.GetCommunity)
	}
}
