package handler

import (
	"strconv"

	"myApp/model"
	"myApp/pkg/response"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

// LandlordHandler 房东处理器结构体，负责处理房东相关的HTTP请求
type LandlordHandler struct {
	service service.LandlordService
}

// NewLandlordHandler 创建房东处理器实例，注入房东服务依赖
func NewLandlordHandler(s service.LandlordService) *LandlordHandler {
	return &LandlordHandler{service: s}
}

// CreateLandlord 创建房东信息
func (h *LandlordHandler) CreateLandlord(c *gin.Context) {
	var landlord model.Landlord
	if err := c.ShouldBindJSON(&landlord); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}
	landlord.UserID = userID.(uint)

	if err := h.service.CreateLandlord(&landlord); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, landlord)
}

// GetLandlordProfile 获取房东个人资料
func (h *LandlordHandler) GetLandlordProfile(c *gin.Context) {
	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	landlord, err := h.service.GetLandlordByUserID(userID.(uint))
	if err != nil {
		response.NotFound(c, "房东信息不存在")
		return
	}

	response.Success(c, landlord)
}

// UpdateLandlord 更新房东信息
func (h *LandlordHandler) UpdateLandlord(c *gin.Context) {
	var landlord model.Landlord
	if err := c.ShouldBindJSON(&landlord); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 检查是否为本人操作
	existingLandlord, err := h.service.GetLandlordByUserID(userID.(uint))
	if err != nil {
		response.NotFound(c, "房东信息不存在")
		return
	}

	// 设置ID
	landlord.ID = existingLandlord.ID

	if err := h.service.UpdateLandlord(&landlord); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, landlord)
}

// VerifyLandlord 管理员验证房东身份
func (h *LandlordHandler) VerifyLandlord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的房东ID")
		return
	}

	// 从上下文获取用户ID和类型（由JWT中间件设置）
	userType, exists := c.Get("userType")
	if !exists || userType.(int) != 2 { // 2-管理员
		response.Forbidden(c, "无权进行此操作")
		return
	}

	if err := h.service.VerifyLandlord(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "房东验证成功"})
}