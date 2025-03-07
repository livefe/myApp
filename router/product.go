package router

import (
	"myApp/handler/product"
	"myApp/repository"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

func InitProductRouter(r *gin.Engine) {
	productRepo := repository.NewProductRepository()
	productCategoryRepo := repository.NewProductCategoryRepository()
	productService := service.NewProductService(productRepo, productCategoryRepo)
	handler := product.NewHandler(productService)
	productGroup := r.Group("/api/product")
	{
		productGroup.POST("/create", handler.CreateProduct)
		productGroup.GET("/list", handler.GetAllProducts)
		productGroup.GET("/:id", handler.GetProduct)
	}
}
