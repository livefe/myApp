package favorite

import (
	"myApp/dto/house"
	"time"
)

// 收藏基本信息DTO
type BasicInfoDTO struct {
	ID        uint      `json:"id"`         // 收藏ID
	UserID    uint      `json:"user_id"`    // 用户ID
	HouseID   uint      `json:"house_id"`   // 房源ID
	Notes     string    `json:"notes"`      // 收藏备注
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// 收藏详细信息DTO
type DetailDTO struct {
	ID        uint               `json:"id"`         // 收藏ID
	UserID    uint               `json:"user_id"`    // 用户ID
	HouseID   uint               `json:"house_id"`   // 房源ID
	Notes     string             `json:"notes"`      // 收藏备注
	House     house.BasicInfoDTO `json:"house"`      // 房源基本信息
	CreatedAt time.Time          `json:"created_at"` // 创建时间
	UpdatedAt time.Time          `json:"updated_at"` // 更新时间
}

// 收藏列表响应DTO
type ListResponse struct {
	Total int         `json:"total"` // 总数
	List  []DetailDTO `json:"list"`  // 列表
}
