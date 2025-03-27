package user

import (
	"time"
)

// 用户基本信息响应DTO
type BasicInfoDTO struct {
	ID       uint   `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名
	Avatar   string `json:"avatar"`   // 头像URL
}

// 用户详细信息响应DTO
type DetailDTO struct {
	ID        uint      `json:"id"`         // 用户ID
	Username  string    `json:"username"`   // 用户名
	Phone     string    `json:"phone"`      // 手机号
	Email     string    `json:"email"`      // 电子邮箱
	RealName  string    `json:"real_name"`  // 真实姓名
	Avatar    string    `json:"avatar"`     // 头像URL
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// 用户登录响应DTO
type LoginResponse struct {
	Token     string    `json:"token"`      // JWT令牌
	ExpiresAt time.Time `json:"expires_at"` // 过期时间
	User      DetailDTO `json:"user"`       // 用户信息
}

// 用户列表响应DTO
type ListResponse struct {
	Total int            `json:"total"` // 总数
	List  []BasicInfoDTO `json:"list"`  // 列表
}
