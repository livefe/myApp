package router

import (
	"myApp/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	// 加载全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.RateLimiter())
	r.Use(middleware.JWTAuth())

	// 初始化子路由
	InitUserRouter(r)
	InitCommunityRouter(r)
	InitOrderRouter(r)
	InitProductRouter(r)
}
