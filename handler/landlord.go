package handler

import (
	"strconv"

	"myApp/dto/landlord"
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
	var req landlord.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 验证请求参数
	if err := landlord.ValidateRegisterRequest(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 将DTO转换为模型
	landlordModel := model.Landlord{
		UserID:       userID.(uint),
		RealName:     req.RealName,
		IDNumber:     req.IDNumber,
		PhoneNumber:  req.PhoneNumber,
		Address:      req.Address,
		IdCardFront:  req.IdCardFront,
		IdCardBack:   req.IdCardBack,
		BankAccount:  req.BankAccount,
		BankName:     req.BankName,
		AccountName:  req.AccountName,
		Introduction: req.Introduction,
		Verified:     false, // 默认未认证
	}

	if err := h.service.CreateLandlord(&landlordModel); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// 将模型转换为DTO
	landlordDTO := landlord.DetailDTO{
		ID:           landlordModel.ID,
		UserID:       landlordModel.UserID,
		RealName:     landlordModel.RealName,
		IDNumber:     landlordModel.IDNumber,
		PhoneNumber:  landlordModel.PhoneNumber,
		Address:      landlordModel.Address,
		Verified:     landlordModel.Verified,
		IdCardFront:  landlordModel.IdCardFront,
		IdCardBack:   landlordModel.IdCardBack,
		BankAccount:  landlordModel.BankAccount,
		BankName:     landlordModel.BankName,
		AccountName:  landlordModel.AccountName,
		Introduction: landlordModel.Introduction,
		Rating:       landlordModel.Rating,
		CreatedAt:    landlordModel.CreatedAt,
	}

	response.Success(c, landlordDTO)
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

	// 解析请求数据
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 更新房东信息，只更新请求中包含的字段
	for key, value := range updateData {
		switch key {
		case "phone_number":
			if phoneNumber, ok := value.(string); ok {
				existingLandlord.PhoneNumber = phoneNumber
			}
		case "address":
			if address, ok := value.(string); ok {
				existingLandlord.Address = address
			}
		case "id_card_front":
			if idCardFront, ok := value.(string); ok {
				existingLandlord.IdCardFront = idCardFront
			}
		case "id_card_back":
			if idCardBack, ok := value.(string); ok {
				existingLandlord.IdCardBack = idCardBack
			}
		case "bank_account":
			if bankAccount, ok := value.(string); ok {
				existingLandlord.BankAccount = bankAccount
			}
		case "bank_name":
			if bankName, ok := value.(string); ok {
				existingLandlord.BankName = bankName
			}
		case "account_name":
			if accountName, ok := value.(string); ok {
				existingLandlord.AccountName = accountName
			}
		case "introduction":
			if introduction, ok := value.(string); ok {
				existingLandlord.Introduction = introduction
			}
		case "real_name":
			if realName, ok := value.(string); ok {
				existingLandlord.RealName = realName
			}
		case "id_number":
			if idNumber, ok := value.(string); ok {
				existingLandlord.IDNumber = idNumber
			}
		}
	}

	if err := h.service.UpdateLandlord(existingLandlord); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, existingLandlord)
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
