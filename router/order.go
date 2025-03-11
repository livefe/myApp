package router

import (
	"myApp/handler"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(r *gin.Engine) {
	orderRepo := repository.NewOrderRepository()
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)
	orderGroup := r.Group("/api/order")
	{
		orderGroup.POST("/create", orderHandler.CreateOrder)
		orderGroup.GET("/:id", orderHandler.GetOrder)
	}
}
