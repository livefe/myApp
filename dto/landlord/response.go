package landlord

import (
	"time"
)

// 房东基本信息DTO
type BasicInfoDTO struct {
	ID          uint      `json:"id"`           // 房东ID
	UserID      uint      `json:"user_id"`      // 关联的用户ID
	RealName    string    `json:"real_name"`    // 真实姓名
	PhoneNumber string    `json:"phone_number"` // 联系电话
	Verified    bool      `json:"verified"`     // 是否已认证
	Rating      float64   `json:"rating"`       // 房东评分
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
}

// 房东详细信息DTO
type DetailDTO struct {
	ID           uint      `json:"id"`            // 房东ID
	UserID       uint      `json:"user_id"`       // 关联的用户ID
	RealName     string    `json:"real_name"`     // 真实姓名
	IDNumber     string    `json:"id_number"`     // 身份证号
	PhoneNumber  string    `json:"phone_number"`  // 联系电话
	Address      string    `json:"address"`       // 联系地址
	Verified     bool      `json:"verified"`      // 是否已认证
	IdCardFront  string    `json:"id_card_front"` // 身份证正面照片URL
	IdCardBack   string    `json:"id_card_back"`  // 身份证背面照片URL
	BankAccount  string    `json:"bank_account"`  // 银行账号
	BankName     string    `json:"bank_name"`     // 开户行名称
	AccountName  string    `json:"account_name"`  // 开户人姓名
	Introduction string    `json:"introduction"`  // 房东自我介绍
	Rating       float64   `json:"rating"`        // 房东评分
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
}

// 房东列表响应DTO
type ListResponse struct {
	Total int            `json:"total"` // 总数
	List  []BasicInfoDTO `json:"list"`  // 列表
}
