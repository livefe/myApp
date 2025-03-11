package handler

import (
	"strconv"

	"myApp/model"
	"myApp/pkg/response"
	"myApp/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(s service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}
	order.UserID = userID.(uint)

	createdOrder, err := h.service.CreateOrder(&order)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, createdOrder)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的订单ID")
		return
	}

	order, err := h.service.GetOrder(uint(id))
	if err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	response.Success(c, order)
}