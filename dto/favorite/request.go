package favorite

import (
	"github.com/go-playground/validator/v10"
)

// 添加收藏请求DTO
type AddRequest struct {
	HouseID uint   `json:"house_id" binding:"required" example:"1"`            // 房源ID
	Notes   string `json:"notes" binding:"omitempty" example:"这套房子采光很好，地段也不错"` // 收藏备注
}

// 收藏查询请求DTO
type QueryRequest struct {
	UserID    uint   `json:"user_id" form:"user_id" example:"1"`          // 用户ID
	HouseID   uint   `json:"house_id" form:"house_id" example:"1"`        // 房源ID
	Keyword   string `json:"keyword" form:"keyword" example:"采光好"`        // 关键词（搜索备注）
	Page      int    `json:"page" form:"page" example:"1"`                // 页码
	PageSize  int    `json:"page_size" form:"page_size" example:"10"`     // 每页数量
	SortBy    string `json:"sort_by" form:"sort_by" example:"created_at"` // 排序字段
	SortOrder string `json:"sort_order" form:"sort_order" example:"desc"` // 排序方向
}

// ValidateAddRequest 验证添加收藏请求
func ValidateAddRequest(req AddRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}
