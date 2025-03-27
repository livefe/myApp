package user

import (
	"github.com/go-playground/validator/v10"
)

// 用户注册请求DTO
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"zhangsan"`            // 用户名
	Password string `json:"password" binding:"required,min=6,max=20" example:"password123"`         // 密码
	Phone    string `json:"phone" binding:"required,len=11" example:"13800138000"`                  // 手机号
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`              // 电子邮箱
	RealName string `json:"real_name" binding:"omitempty" example:"张三"`                             // 真实姓名
	IdCard   string `json:"id_card" binding:"omitempty,len=18" example:"110101199001011234"`        // 身份证号
	Avatar   string `json:"avatar" binding:"omitempty,url" example:"http://example.com/avatar.jpg"` // 头像URL
}

// 用户登录请求DTO
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"zhangsan"`    // 用户名
	Password string `json:"password" binding:"required" example:"password123"` // 密码
}

// 用户信息更新请求DTO
type UpdateRequest struct {
	Phone    string `json:"phone" binding:"omitempty,len=11" example:"13800138000"`                 // 手机号
	Email    string `json:"email" binding:"omitempty,email" example:"user@example.com"`             // 电子邮箱
	RealName string `json:"real_name" binding:"omitempty" example:"张三"`                             // 真实姓名
	IdCard   string `json:"id_card" binding:"omitempty,len=18" example:"110101199001011234"`        // 身份证号
	Avatar   string `json:"avatar" binding:"omitempty,url" example:"http://example.com/avatar.jpg"` // 头像URL
}

// 密码修改请求DTO
type PasswordUpdateRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"oldpassword123"`              // 旧密码
	NewPassword string `json:"new_password" binding:"required,min=6,max=20" example:"newpassword123"` // 新密码
}

// ValidateRegisterRequest 验证用户注册请求
func ValidateRegisterRequest(req RegisterRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}

// ValidateLoginRequest 验证用户登录请求
func ValidateLoginRequest(req LoginRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}

// ValidateUpdateRequest 验证用户信息更新请求
func ValidateUpdateRequest(req UpdateRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}
