package handler

import (
	"myApp/config"
	"myApp/model"
	"myApp/pkg/response"
	"myApp/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// UserHandler 用户处理器结构体，负责处理用户相关的HTTP请求
type UserHandler struct {
	service service.UserService
}

// NewUserHandler 创建用户处理器实例，注入用户服务依赖
func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Register 用户注册处理函数
func (h *UserHandler) Register(c *gin.Context) {
	// 绑定并验证请求参数
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务层进行用户注册
	createdUser, err := h.service.Register(&user)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// 返回成功响应
	response.Success(c, createdUser)
}

// Login 用户登录处理函数
func (h *UserHandler) Login(c *gin.Context) {
	// 绑定并验证登录凭证
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 调用服务层进行用户登录验证
	user, err := h.service.Login(credentials.Username, credentials.Password)
	if err != nil {
		response.Unauthorized(c, "无效的凭证")
		return
	}

	// 生成JWT令牌
	claims := jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Duration(config.Conf.JWT.Expire) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Conf.JWT.Secret))
	if err != nil {
		response.ServerError(c, "生成令牌失败")
		return
	}

	// 返回成功响应，包含令牌和用户信息
	response.SuccessWithToken(c, tokenString, user)
}

// GetUserInfo 获取用户信息处理函数
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	// 调用服务层获取用户信息
	user, err := h.service.GetUserProfile(userID.(uint))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 返回成功响应
	response.Success(c, user)
}