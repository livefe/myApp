package router

import (
	"myApp/handler/user"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.Engine) {
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	handler := user.NewUserHandler(userService)
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/register", handler.Register)
		userGroup.POST("/login", handler.Login)
		userGroup.GET("/info", handler.GetUserInfo)
	}
}
