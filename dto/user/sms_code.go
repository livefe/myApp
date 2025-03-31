package user

import (
	"github.com/go-playground/validator/v10"
)

// 发送短信验证码请求DTO
type SendSMSCodeRequest struct {
	Phone string `json:"phone" binding:"required,len=11" example:"13800138000"` // 手机号
}

// 短信验证码登录请求DTO
type SMSCodeLoginRequest struct {
	Phone string `json:"phone" binding:"required,len=11" example:"13800138000"` // 手机号
	Code  string `json:"code" binding:"required,len=6" example:"123456"`        // 验证码
}

// ValidateSendSMSCodeRequest 验证发送短信验证码请求
func ValidateSendSMSCodeRequest(req SendSMSCodeRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}

// ValidateSMSCodeLoginRequest 验证短信验证码登录请求
func ValidateSMSCodeLoginRequest(req SMSCodeLoginRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}
