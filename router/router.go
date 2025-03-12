package router

import (
	"myApp/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置所有路由和中间件
func SetupRouter(r *gin.Engine) {

	// 加载全局中间件
	r.Use(middleware.CORS())      // 跨域资源共享中间件
	r.Use(middleware.Logger())    // 日志记录中间件
	r.Use(middleware.RateLimiter()) // 请求速率限制中间件
	r.Use(middleware.JWTAuth())   // JWT认证中间件

	// 初始化子路由
	InitUserRouter(r)     // 初始化用户相关路由
	InitHouseRouter(r)    // 初始化房源相关路由
	InitViewingRouter(r)  // 初始化预约看房相关路由
	InitFavoriteRouter(r) // 初始化收藏相关路由
	RegisterLandlordRoutes(r.Group("/api")) // 初始化房东相关路由
}
