package handler

import (
	"strconv"

	"myApp/model"
	"myApp/pkg/response"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{service: productService}
}

// CreateProduct
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}
	product.CreatorID = userID.(uint)

	if err := h.service.CreateProduct(&product); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, product)
}

// GetProduct
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID格式")
		return
	}

	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		response.NotFound(c, "产品不存在")
		return
	}

	response.Success(c, product)
}

// GetAllProducts
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	// 创建一个空的参数映射，可以根据请求参数进行填充
	params := make(map[string]interface{})
	
	// 可以从查询参数中获取过滤条件
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			params["category_id"] = uint(categoryID)
		}
	}

	if communityIDStr := c.Query("community_id"); communityIDStr != "" {
		communityID, err := strconv.ParseUint(communityIDStr, 10, 32)
		if err == nil {
			params["community_id"] = uint(communityID)
		}
	}

	products, err := h.service.GetAllProducts(params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, products)
}