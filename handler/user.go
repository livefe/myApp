package handler

import (
	"myApp/config"
	"myApp/dto/user"
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
	var req user.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 验证请求参数
	if err := user.ValidateRegisterRequest(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 将DTO转换为模型
	userModel := model.User{
		Username: req.Username,
		Password: req.Password,
		Phone:    req.Phone,
		Email:    req.Email,
		RealName: req.RealName,
		IdCard:   req.IdCard,
		Avatar:   req.Avatar,
	}

	// 调用服务层进行用户注册
	createdUser, err := h.service.Register(&userModel)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// 将模型转换为DTO
	userDTO := user.DetailDTO{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Phone:     createdUser.Phone,
		Email:     createdUser.Email,
		RealName:  createdUser.RealName,
		Avatar:    createdUser.Avatar,
		CreatedAt: createdUser.CreatedAt,
	}

	// 返回成功响应
	response.Success(c, userDTO)
}

// Login 用户登录处理函数
func (h *UserHandler) Login(c *gin.Context) {
	// 绑定并验证登录凭证
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 验证请求参数
	if err := user.ValidateLoginRequest(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务层进行用户登录验证
	userModel, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, "无效的凭证")
		return
	}

	// 生成JWT令牌
	claims := jwt.MapClaims{
		"userID": userModel.ID,
		"exp":    time.Now().Add(time.Duration(config.Conf.JWT.Expire) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Conf.JWT.Secret))
	if err != nil {
		response.ServerError(c, "生成令牌失败")
		return
	}

	// 将模型转换为DTO
	userDTO := user.DetailDTO{
		ID:        userModel.ID,
		Username:  userModel.Username,
		Phone:     userModel.Phone,
		Email:     userModel.Email,
		RealName:  userModel.RealName,
		Avatar:    userModel.Avatar,
		CreatedAt: userModel.CreatedAt,
	}

	// 创建登录响应DTO
	loginResponse := user.LoginResponse{
		Token:     tokenString,
		ExpiresAt: time.Now().Add(time.Duration(config.Conf.JWT.Expire) * time.Second),
		User:      userDTO,
	}

	// 返回成功响应
	response.Success(c, loginResponse)
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
	userModel, err := h.service.GetUserProfile(userID.(uint))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 将模型转换为DTO
	userDTO := user.DetailDTO{
		ID:        userModel.ID,
		Username:  userModel.Username,
		Phone:     userModel.Phone,
		Email:     userModel.Email,
		RealName:  userModel.RealName,
		Avatar:    userModel.Avatar,
		CreatedAt: userModel.CreatedAt,
	}

	// 返回成功响应
	response.Success(c, userDTO)
}
