package handler

import (
	"myApp/config"
	"myApp/dto/user"
	"myApp/pkg/response"
	"myApp/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SMSCodeHandler 短信验证码处理器结构体
type SMSCodeHandler struct {
	smsCodeService service.SMSCodeService
}

// NewSMSCodeHandler 创建短信验证码处理器实例
func NewSMSCodeHandler(s service.SMSCodeService) *SMSCodeHandler {
	return &SMSCodeHandler{smsCodeService: s}
}

// SendCode 发送短信验证码处理函数
func (h *SMSCodeHandler) SendCode(c *gin.Context) {
	// 绑定并验证请求参数
	var req user.SendSMSCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 验证请求参数
	if err := user.ValidateSendSMSCodeRequest(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务层发送验证码
	success, err := h.smsCodeService.SendCode(req.Phone)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	if !success {
		response.ServerError(c, "发送验证码失败")
		return
	}

	// 返回成功响应
	response.Success(c, gin.H{"message": "验证码已发送"})
}

// LoginByCode 短信验证码登录处理函数
func (h *SMSCodeHandler) LoginByCode(c *gin.Context) {
	// 绑定并验证请求参数
	var req user.SMSCodeLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 验证请求参数
	if err := user.ValidateSMSCodeLoginRequest(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务层验证码登录
	userModel, err := h.smsCodeService.LoginByCode(req.Phone, req.Code)
	if err != nil {
		response.Unauthorized(c, err.Error())
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

	// 返回成功响应
	response.Success(c, gin.H{
		"token": tokenString,
		"user":  userDTO,
	})
}
