package landlord

import (
	"github.com/go-playground/validator/v10"
)

// 房东注册请求DTO
type RegisterRequest struct {
	UserID       uint   `json:"user_id" binding:"required" example:"1"`                                      // 关联的用户ID
	RealName     string `json:"real_name" binding:"required" example:"张三"`                                   // 真实姓名
	IDNumber     string `json:"id_number" binding:"required,len=18" example:"110101199001011234"`            // 身份证号
	PhoneNumber  string `json:"phone_number" binding:"required,len=11" example:"13800138000"`                // 联系电话
	Address      string `json:"address" binding:"required" example:"北京市朝阳区建国路1号"`                            // 联系地址
	IdCardFront  string `json:"id_card_front" binding:"required,url" example:"http://example.com/front.jpg"` // 身份证正面照片URL
	IdCardBack   string `json:"id_card_back" binding:"required,url" example:"http://example.com/back.jpg"`   // 身份证背面照片URL
	BankAccount  string `json:"bank_account" binding:"required" example:"6222021234567890123"`               // 银行账号
	BankName     string `json:"bank_name" binding:"required" example:"中国工商银行"`                               // 开户行名称
	AccountName  string `json:"account_name" binding:"required" example:"张三"`                                // 开户人姓名
	Introduction string `json:"introduction" binding:"omitempty" example:"专业的房产经纪人，有多年租赁经验"`                 // 房东自我介绍
}

// 房东信息更新请求DTO
type UpdateRequest struct {
	RealName     string `json:"real_name" binding:"omitempty" example:"张三"`                                   // 真实姓名
	PhoneNumber  string `json:"phone_number" binding:"omitempty,len=11" example:"13800138000"`                // 联系电话
	Address      string `json:"address" binding:"omitempty" example:"北京市朝阳区建国路1号"`                            // 联系地址
	IdCardFront  string `json:"id_card_front" binding:"omitempty,url" example:"http://example.com/front.jpg"` // 身份证正面照片URL
	IdCardBack   string `json:"id_card_back" binding:"omitempty,url" example:"http://example.com/back.jpg"`   // 身份证背面照片URL
	BankAccount  string `json:"bank_account" binding:"omitempty" example:"6222021234567890123"`               // 银行账号
	BankName     string `json:"bank_name" binding:"omitempty" example:"中国工商银行"`                               // 开户行名称
	AccountName  string `json:"account_name" binding:"omitempty" example:"张三"`                                // 开户人姓名
	Introduction string `json:"introduction" binding:"omitempty" example:"专业的房产经纪人，有多年租赁经验"`                  // 房东自我介绍
	Verified     bool   `json:"verified" example:"true"`                                                      // 是否已认证
}

// 房东查询请求DTO
type QueryRequest struct {
	Keyword   string  `json:"keyword" form:"keyword" example:"张三"`         // 关键词
	Verified  bool    `json:"verified" form:"verified" example:"true"`     // 是否已认证
	MinRating float64 `json:"min_rating" form:"min_rating" example:"4.5"`  // 最低评分
	UserID    uint    `json:"user_id" form:"user_id" example:"1"`          // 关联的用户ID
	Page      int     `json:"page" form:"page" example:"1"`                // 页码
	PageSize  int     `json:"page_size" form:"page_size" example:"10"`     // 每页数量
	SortBy    string  `json:"sort_by" form:"sort_by" example:"created_at"` // 排序字段
	SortOrder string  `json:"sort_order" form:"sort_order" example:"desc"` // 排序方向
}

// ValidateRegisterRequest 验证房东注册请求
func ValidateRegisterRequest(req RegisterRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}

// ValidateUpdateRequest 验证房东信息更新请求
func ValidateUpdateRequest(req UpdateRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}
