package handler

import (
	"strconv"

	"myApp/dto/house"
	"myApp/model"
	"myApp/pkg/response"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

// HouseHandler 房源处理器结构体，负责处理房源相关的HTTP请求
type HouseHandler struct {
	service service.HouseService
}

// NewHouseHandler 创建房源处理器实例，注入房源服务依赖
func NewHouseHandler(s service.HouseService) *HouseHandler {
	return &HouseHandler{service: s}
}

// CreateHouse 创建房源
func (h *HouseHandler) CreateHouse(c *gin.Context) {
	var req house.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 验证请求参数
	if err := house.ValidateCreateRequest(req); err != nil {
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
	houseModel := model.House{
		Title:       req.Title,
		Description: req.Description,
		Address:     req.Address,
		Area:        req.Area,
		Floor:       req.Floor,
		TotalFloor:  req.TotalFloor,
		Rooms:       req.Rooms,
		Halls:       req.Halls,
		Bathrooms:   req.Bathrooms,
		RentPrice:   req.RentPrice,
		Deposit:     req.Deposit,
		PaymentType: req.PaymentType,
		HouseType:   req.HouseType,
		Orientation: req.Orientation,
		Decoration:  req.Decoration,
		Facilities:  req.Facilities,
		Images:      req.Images,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		IsElevator:  req.IsElevator,
		LandlordID:  userID.(uint),
		Status:      1, // 默认上架状态
	}

	if err := h.service.CreateHouse(&houseModel); err != nil {
		response.ServerError(c, "创建房源失败")
		return
	}

	// 将模型转换为DTO
	houseDTO := house.DetailDTO{
		ID:          houseModel.ID,
		Title:       houseModel.Title,
		Description: houseModel.Description,
		Address:     houseModel.Address,
		Area:        houseModel.Area,
		Floor:       houseModel.Floor,
		TotalFloor:  houseModel.TotalFloor,
		Rooms:       houseModel.Rooms,
		Halls:       houseModel.Halls,
		Bathrooms:   houseModel.Bathrooms,
		RentPrice:   houseModel.RentPrice,
		Deposit:     houseModel.Deposit,
		PaymentType: houseModel.PaymentType,
		HouseType:   houseModel.HouseType,
		Orientation: houseModel.Orientation,
		Decoration:  houseModel.Decoration,
		Facilities:  houseModel.Facilities,
		Images:      houseModel.Images,
		Latitude:    houseModel.Latitude,
		Longitude:   houseModel.Longitude,
		IsElevator:  houseModel.IsElevator,
		Status:      houseModel.Status,
		LandlordID:  houseModel.LandlordID,
		ViewCount:   houseModel.ViewCount,
		CreatedAt:   houseModel.CreatedAt,
		UpdatedAt:   houseModel.UpdatedAt,
	}

	response.Success(c, houseDTO)
}

// GetHouse 获取房源详情
func (h *HouseHandler) GetHouse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的房源ID")
		return
	}

	// 增加浏览次数
	h.service.IncrementViewCount(uint(id))

	houseModel, err := h.service.GetHouseByID(uint(id))
	if err != nil {
		response.NotFound(c, "房源不存在")
		return
	}

	// 将模型转换为DTO
	houseDTO := house.DetailDTO{
		ID:          houseModel.ID,
		Title:       houseModel.Title,
		Description: houseModel.Description,
		Address:     houseModel.Address,
		Area:        houseModel.Area,
		Floor:       houseModel.Floor,
		TotalFloor:  houseModel.TotalFloor,
		Rooms:       houseModel.Rooms,
		Halls:       houseModel.Halls,
		Bathrooms:   houseModel.Bathrooms,
		RentPrice:   houseModel.RentPrice,
		Deposit:     houseModel.Deposit,
		PaymentType: houseModel.PaymentType,
		HouseType:   houseModel.HouseType,
		Orientation: houseModel.Orientation,
		Decoration:  houseModel.Decoration,
		Facilities:  houseModel.Facilities,
		Images:      houseModel.Images,
		Latitude:    houseModel.Latitude,
		Longitude:   houseModel.Longitude,
		IsElevator:  houseModel.IsElevator,
		Status:      houseModel.Status,
		LandlordID:  houseModel.LandlordID,
		ViewCount:   houseModel.ViewCount,
		CreatedAt:   houseModel.CreatedAt,
		UpdatedAt:   houseModel.UpdatedAt,
	}

	response.Success(c, houseDTO)
}

// GetAllHouses 获取房源列表
func (h *HouseHandler) GetAllHouses(c *gin.Context) {
	// 解析查询参数
	params := make(map[string]interface{})

	// 状态筛选
	if statusStr := c.Query("status"); statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err == nil {
			params["status"] = status
		}
	}

	// 房东ID筛选
	if landlordIDStr := c.Query("landlord_id"); landlordIDStr != "" {
		landlordID, err := strconv.ParseUint(landlordIDStr, 10, 64)
		if err == nil {
			params["landlord_id"] = uint(landlordID)
		}
	}

	// 价格范围筛选
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err == nil {
			params["min_price"] = minPrice
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err == nil {
			params["max_price"] = maxPrice
		}
	}

	// 房间数筛选
	if roomsStr := c.Query("rooms"); roomsStr != "" {
		rooms, err := strconv.Atoi(roomsStr)
		if err == nil {
			params["rooms"] = rooms
		}
	}

	// 房屋类型筛选
	if houseTypeStr := c.Query("house_type"); houseTypeStr != "" {
		houseType, err := strconv.Atoi(houseTypeStr)
		if err == nil {
			params["house_type"] = houseType
		}
	}

	// 关键词搜索
	if keyword := c.Query("keyword"); keyword != "" {
		params["keyword"] = keyword
	}

	// 排序
	if orderBy := c.Query("order_by"); orderBy != "" {
		params["order_by"] = orderBy
	}

	// 分页
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			params["limit"] = limit

			if offsetStr := c.Query("offset"); offsetStr != "" {
				offset, err := strconv.Atoi(offsetStr)
				if err == nil && offset >= 0 {
					params["offset"] = offset
				}
			}
		}
	}

	houses, err := h.service.GetAllHouses(params)
	if err != nil {
		response.ServerError(c, "获取房源列表失败")
		return
	}

	response.Success(c, houses)
}

// UpdateHouse 更新房源
func (h *HouseHandler) UpdateHouse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的房源ID")
		return
	}

	// 检查房源是否存在
	existingHouse, err := h.service.GetHouseByID(uint(id))
	if err != nil {
		response.NotFound(c, "房源不存在")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 检查是否为房东本人
	if existingHouse.LandlordID != userID.(uint) {
		response.Forbidden(c, "无权修改该房源")
		return
	}

	var house model.House
	if err := c.ShouldBindJSON(&house); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 设置ID和房东ID
	house.ID = uint(id)
	house.LandlordID = userID.(uint)

	if err := h.service.UpdateHouse(&house); err != nil {
		response.ServerError(c, "更新房源失败")
		return
	}

	response.Success(c, house)
}

// DeleteHouse 删除房源
func (h *HouseHandler) DeleteHouse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的房源ID")
		return
	}

	// 检查房源是否存在
	existingHouse, err := h.service.GetHouseByID(uint(id))
	if err != nil {
		response.NotFound(c, "房源不存在")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 检查是否为房东本人
	if existingHouse.LandlordID != userID.(uint) {
		response.Forbidden(c, "无权删除该房源")
		return
	}

	if err := h.service.DeleteHouse(uint(id)); err != nil {
		response.ServerError(c, "删除房源失败")
		return
	}

	response.Success(c, nil)
}

// GetLandlordHouses 获取房东的所有房源
func (h *HouseHandler) GetLandlordHouses(c *gin.Context) {
	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	houses, err := h.service.GetHousesByLandlordID(userID.(uint))
	if err != nil {
		response.ServerError(c, "获取房源列表失败")
		return
	}

	response.Success(c, houses)
}
