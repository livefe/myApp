package handler

import (
	"strconv"

	"myApp/model"
	"myApp/pkg/response"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

// FavoriteHandler 收藏处理器结构体，负责处理收藏相关的HTTP请求
type FavoriteHandler struct {
	service service.FavoriteService
}

// NewFavoriteHandler 创建收藏处理器实例，注入收藏服务依赖
func NewFavoriteHandler(s service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: s}
}

// AddFavorite 添加收藏
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	var favorite model.Favorite
	if err := c.ShouldBindJSON(&favorite); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}
	favorite.UserID = userID.(uint)

	if err := h.service.AddFavorite(&favorite); err != nil {
		response.ServerError(c, "添加收藏失败")
		return
	}

	response.Success(c, favorite)
}

// RemoveFavorite 删除收藏
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的收藏ID")
		return
	}

	// 检查收藏是否存在
	favorite, err := h.service.GetFavoriteByID(uint(id))
	if err != nil {
		response.NotFound(c, "收藏记录不存在")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 检查是否为用户本人的收藏
	if favorite.UserID != userID.(uint) {
		response.Forbidden(c, "无权删除该收藏")
		return
	}

	if err := h.service.RemoveFavorite(uint(id)); err != nil {
		response.ServerError(c, "删除收藏失败")
		return
	}

	response.Success(c, nil)
}

// GetUserFavorites 获取用户的所有收藏
func (h *FavoriteHandler) GetUserFavorites(c *gin.Context) {
	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	favorites, err := h.service.GetUserFavorites(userID.(uint))
	if err != nil {
		response.ServerError(c, "获取收藏列表失败")
		return
	}

	response.Success(c, favorites)
}

// ToggleFavorite 切换收藏状态
func (h *FavoriteHandler) ToggleFavorite(c *gin.Context) {
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

	// 获取备注信息
	var data struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		data.Notes = ""
	}

	if err := h.service.ToggleFavorite(userID.(uint), uint(houseID), data.Notes); err != nil {
		response.ServerError(c, "操作收藏失败")
		return
	}

	// 检查当前状态
	isFav, err := h.service.IsFavorite(userID.(uint), uint(houseID))
	if err != nil {
		response.ServerError(c, "获取收藏状态失败")
		return
	}

	response.Success(c, gin.H{"is_favorite": isFav})
}

// CheckFavorite 检查是否已收藏
func (h *FavoriteHandler) CheckFavorite(c *gin.Context) {
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

	isFav, err := h.service.IsFavorite(userID.(uint), uint(houseID))
	if err != nil {
		response.ServerError(c, "获取收藏状态失败")
		return
	}

	response.Success(c, gin.H{"is_favorite": isFav})
}