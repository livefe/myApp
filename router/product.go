package router

import (
	"myApp/handler"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

func InitProductRouter(r *gin.Engine) {
	productRepo := repository.NewProductRepository()
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)
	productGroup := r.Group("/api/product")
	{
		productGroup.POST("/create", productHandler.CreateProduct)
		productGroup.GET("/list", productHandler.GetAllProducts)
		productGroup.GET("/:id", productHandler.GetProduct)
	}
}
