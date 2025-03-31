package handler

import (
	"strconv"
	"time"

	"myApp/dto/viewing"
	"myApp/model"
	"myApp/pkg/response"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

// ViewingHandler 预约看房处理器结构体，负责处理预约看房相关的HTTP请求
type ViewingHandler struct {
	service      service.ViewingService
	houseService service.HouseService
}

// NewViewingHandler 创建预约看房处理器实例，注入预约看房服务依赖
func NewViewingHandler(s service.ViewingService, hs service.HouseService) *ViewingHandler {
	return &ViewingHandler{service: s, houseService: hs}
}

// CreateViewing 创建预约看房
func (h *ViewingHandler) CreateViewing(c *gin.Context) {
	// 绑定并验证请求参数
	var req viewing.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 验证请求参数
	if err := viewing.ValidateCreateRequest(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 验证预约时间是否合法（不能是过去的时间）
	if req.ViewDate.Before(time.Now()) {
		response.BadRequest(c, "预约时间不能是过去的时间")
		return
	}

	// 将DTO转换为模型
	viewingModel := model.Viewing{
		HouseID:      req.HouseID,
		UserID:       userID.(uint),
		ViewingTime:  req.ViewDate,
		Remark:       req.Message,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       0, // 默认待确认状态
	}

	if err := h.service.CreateViewing(&viewingModel); err != nil {
		response.ServerError(c, "创建预约看房失败")
		return
	}

	// 返回成功响应，包含预约ID
	response.Success(c, gin.H{"id": viewingModel.ID})
}

// GetViewing 获取预约看房详情
func (h *ViewingHandler) GetViewing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的预约ID")
		return
	}

	viewingModel, err := h.service.GetViewingByID(uint(id))
	if err != nil {
		response.NotFound(c, "预约记录不存在")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 检查是否为预约用户本人或房东
	if viewingModel.UserID != userID.(uint) {
		// 获取房源信息，检查当前用户是否为房东
		house, err := h.houseService.GetHouseByID(viewingModel.HouseID)
		if err != nil || house.LandlordID != userID.(uint) {
			response.Forbidden(c, "无权查看该预约记录")
			return
		}
	}

	// 使用DTO响应结构体构建响应数据
	viewingResp := viewing.ViewingResponse{
		ID:           viewingModel.ID,
		HouseID:      viewingModel.HouseID,
		UserID:       viewingModel.UserID,
		ViewingTime:  viewingModel.ViewingTime,
		Status:       viewingModel.Status,
		StatusText:   viewing.GetStatusText(viewingModel.Status),
		Remark:       viewingModel.Remark,
		ContactName:  viewingModel.ContactName,
		ContactPhone: viewingModel.ContactPhone,
		ConfirmTime:  viewingModel.ConfirmTime,
		CancelTime:   viewingModel.CancelTime,
		CancelReason: viewingModel.CancelReason,
		CreatedAt:    viewingModel.CreatedAt,
	}

	response.Success(c, viewingResp)
}

// GetUserViewings 获取用户的所有预约看房记录
func (h *ViewingHandler) GetUserViewings(c *gin.Context) {
	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	viewings, err := h.service.GetViewingsByUserID(userID.(uint))
	if err != nil {
		response.ServerError(c, "获取预约记录失败")
		return
	}

	// 转换为响应DTO
	viewingResps := make([]viewing.ViewingResponse, 0, len(viewings))
	for _, v := range viewings {
		viewingResps = append(viewingResps, viewing.ViewingResponse{
			ID:           v.ID,
			HouseID:      v.HouseID,
			UserID:       v.UserID,
			ViewingTime:  v.ViewingTime,
			Status:       v.Status,
			StatusText:   viewing.GetStatusText(v.Status),
			Remark:       v.Remark,
			ContactName:  v.ContactName,
			ContactPhone: v.ContactPhone,
			ConfirmTime:  v.ConfirmTime,
			CancelTime:   v.CancelTime,
			CancelReason: v.CancelReason,
			CreatedAt:    v.CreatedAt,
		})
	}

	response.Success(c, gin.H{"viewings": viewingResps})
}

// GetHouseViewings 获取房源的所有预约看房记录
func (h *ViewingHandler) GetHouseViewings(c *gin.Context) {
	houseIDStr := c.Param("house_id")
	houseID, err := strconv.ParseUint(houseIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的房源ID")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 检查当前用户是否为房东
	house, err := h.houseService.GetHouseByID(uint(houseID))
	if err != nil || house.LandlordID != userID.(uint) {
		response.Forbidden(c, "无权查看该房源的预约记录")
		return
	}

	viewings, err := h.service.GetViewingsByHouseID(uint(houseID))
	if err != nil {
		response.ServerError(c, "获取预约记录失败")
		return
	}

	// 转换为响应DTO
	viewingResps := make([]viewing.ViewingResponse, 0, len(viewings))
	for _, v := range viewings {
		viewingResps = append(viewingResps, viewing.ViewingResponse{
			ID:           v.ID,
			HouseID:      v.HouseID,
			UserID:       v.UserID,
			ViewingTime:  v.ViewingTime,
			Status:       v.Status,
			StatusText:   viewing.GetStatusText(v.Status),
			Remark:       v.Remark,
			ContactName:  v.ContactName,
			ContactPhone: v.ContactPhone,
			ConfirmTime:  v.ConfirmTime,
			CancelTime:   v.CancelTime,
			CancelReason: v.CancelReason,
			CreatedAt:    v.CreatedAt,
		})
	}

	response.Success(c, gin.H{"viewings": viewingResps})
}

// ConfirmViewing 确认预约看房
func (h *ViewingHandler) ConfirmViewing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的预约ID")
		return
	}

	// 获取预约记录
	viewingModel, err := h.service.GetViewingByID(uint(id))
	if err != nil {
		response.NotFound(c, "预约记录不存在")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 检查当前用户是否为房东
	house, err := h.houseService.GetHouseByID(viewingModel.HouseID)
	if err != nil || house.LandlordID != userID.(uint) {
		response.Forbidden(c, "无权确认该预约")
		return
	}

	// 确认预约
	if err := h.service.ConfirmViewing(uint(id)); err != nil {
		response.ServerError(c, "确认预约失败")
		return
	}

	response.Success(c, nil)
}

// CompleteViewing 完成预约看房
func (h *ViewingHandler) CompleteViewing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的预约ID")
		return
	}

	// 获取预约记录
	viewingModel, err := h.service.GetViewingByID(uint(id))
	if err != nil {
		response.NotFound(c, "预约记录不存在")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 检查当前用户是否为房东
	house, err := h.houseService.GetHouseByID(viewingModel.HouseID)
	if err != nil || house.LandlordID != userID.(uint) {
		response.Forbidden(c, "无权完成该预约")
		return
	}

	// 完成预约
	if err := h.service.CompleteViewing(uint(id)); err != nil {
		response.ServerError(c, "完成预约失败")
		return
	}

	response.Success(c, nil)
}

// CancelViewing 取消预约看房
func (h *ViewingHandler) CancelViewing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的预约ID")
		return
	}

	// 获取预约记录
	viewingModel, err := h.service.GetViewingByID(uint(id))
	if err != nil {
		response.NotFound(c, "预约记录不存在")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 检查当前用户是否为预约用户本人或房东
	isLandlord := false
	house, err := h.houseService.GetHouseByID(viewingModel.HouseID)
	if err == nil && house.LandlordID == userID.(uint) {
		isLandlord = true
	}

	if viewingModel.UserID != userID.(uint) && !isLandlord {
		response.Forbidden(c, "无权取消该预约")
		return
	}

	// 获取取消原因
	var cancelData struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&cancelData); err != nil {
		cancelData.Reason = "用户取消"
	}

	// 取消预约
	if err := h.service.CancelViewing(uint(id), cancelData.Reason); err != nil {
		response.ServerError(c, "取消预约失败")
		return
	}

	response.Success(c, nil)
}
