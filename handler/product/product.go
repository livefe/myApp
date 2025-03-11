package product

import (
	"strconv"

	"myApp/model"
	"myApp/pkg/response"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.ProductService
}

func NewHandler(productService service.ProductService) *Handler {
	return &Handler{service: productService}
}

// CreateProduct
func (h *Handler) CreateProduct(c *gin.Context) {
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
func (h *Handler) GetProduct(c *gin.Context) {
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
func (h *Handler) GetAllProducts(c *gin.Context) {
	params := make(map[string]interface{})

	if communityIDStr := c.Query("community_id"); communityIDStr != "" {
		communityID, err := strconv.ParseUint(communityIDStr, 10, 32)
		if err == nil {
			params["community_id"] = uint(communityID)
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			params["category_id"] = uint(categoryID)
		}
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err == nil {
			params["status"] = status
		}
	}

	products, err := h.service.GetAllProducts(params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, products)
}

// UpdateProduct
func (h *Handler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID格式")
		return
	}

	existingProduct, err := h.service.GetProductByID(uint(id))
	if err != nil {
		response.NotFound(c, "产品不存在")
		return
	}

	// 检查用户是否为产品创建者
	userID, exists := c.Get("userID")
	if !exists || existingProduct.CreatorID != userID.(uint) {
		response.Fail(c, 403, "无权修改此产品")
		return
	}

	var updatedProduct model.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	updatedProduct.ID = uint(id)
	updatedProduct.CreatorID = existingProduct.CreatorID

	if err := h.service.UpdateProduct(&updatedProduct); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, updatedProduct)
}

// DeleteProduct
func (h *Handler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID格式")
		return
	}

	existingProduct, err := h.service.GetProductByID(uint(id))
	if err != nil {
		response.NotFound(c, "产品不存在")
		return
	}

	// 检查用户是否为产品创建者
	userID, exists := c.Get("userID")
	if !exists || existingProduct.CreatorID != userID.(uint) {
		response.Fail(c, 403, "无权删除此产品")
		return
	}

	if err := h.service.DeleteProduct(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "产品删除成功"})
}

// Category Handlers

// CreateCategory
func (h *Handler) CreateCategory(c *gin.Context) {
	var category model.ProductCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	if err := h.service.CreateCategory(&category); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, category)
}

// GetAllCategories
func (h *Handler) GetAllCategories(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, categories)
}

// GetCategory
func (h *Handler) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID格式")
		return
	}

	category, err := h.service.GetCategoryByID(uint(id))
	if err != nil {
		response.NotFound(c, "分类不存在")
		return
	}

	response.Success(c, category)
}

// UpdateCategory
func (h *Handler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID格式")
		return
	}

	existingCategory, err := h.service.GetCategoryByID(uint(id))
	if err != nil {
		response.NotFound(c, "分类不存在")
		return
	}

	var updatedCategory model.ProductCategory
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	updatedCategory.ID = existingCategory.ID

	if err := h.service.UpdateCategory(&updatedCategory); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, updatedCategory)
}

// DeleteCategory
func (h *Handler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID格式")
		return
	}

	_, err = h.service.GetCategoryByID(uint(id))
	if err != nil {
		response.NotFound(c, "分类不存在")
		return
	}

	// 检查分类是否被使用
	products, err := h.service.GetAllProducts(map[string]interface{}{"category_id": uint(id)})
	if err == nil && len(products) > 0 {
		response.BadRequest(c, "该分类下存在产品，无法删除")
		return
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "分类删除成功"})
}
