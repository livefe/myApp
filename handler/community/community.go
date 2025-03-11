package community

import (
	"myApp/model"
	"myApp/pkg/response"
	"myApp/service"
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
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}
	community.CreatorID = userID.(uint)

	createdCommunity, err := h.service.CreateCommunity(&community)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, createdCommunity)
}

func (h *Handler) GetCommunityList(c *gin.Context) {
	communities, err := h.service.GetAllCommunities()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, communities)
}

func (h *Handler) GetCommunity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID格式")
		return
	}

	community, err := h.service.GetCommunityByID(uint(id))
	if err != nil {
		response.NotFound(c, "社区不存在")
		return
	}

	response.Success(c, community)
}
