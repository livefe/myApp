package router

import (
	"myApp/handler/order"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(r *gin.Engine) {
	orderRepo := repository.NewOrderRepository()
	orderService := service.NewOrderService(orderRepo)
	handler := order.NewHandler(orderService)
	orderGroup := r.Group("/api/order")
	{
		orderGroup.POST("/create", handler.CreateOrder)
		orderGroup.GET("/:id", handler.GetOrder)
	}
}
