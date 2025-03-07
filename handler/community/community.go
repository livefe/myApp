package community

import (
	"myApp/model"
	"myApp/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.CommunityService
}

func NewCommunityHandler(s service.CommunityService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateCommunity(c *gin.Context) {
	var community model.Community
	if err := c.ShouldBindJSON(&community); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	community.CreatorID = userID.(uint)

	createdCommunity, err := h.service.CreateCommunity(&community)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCommunity)
}

func (h *Handler) GetCommunityList(c *gin.Context) {
	communities, err := h.service.GetAllCommunities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, communities)
}

func (h *Handler) GetCommunity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID格式"})
		return
	}

	community, err := h.service.GetCommunityByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "社区不存在"})
		return
	}

	c.JSON(http.StatusOK, community)
}
