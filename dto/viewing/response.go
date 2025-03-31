package viewing

import (
	"myApp/dto/common"
	"time"
)

// ViewingResponse 预约看房响应DTO
type ViewingResponse struct {
	ID           uint       `json:"id" example:"1"`                                        // 预约ID
	HouseID      uint       `json:"house_id" example:"1"`                                  // 房源ID
	UserID       uint       `json:"user_id" example:"1"`                                   // 用户ID
	ViewingTime  time.Time  `json:"viewing_time" example:"2023-07-01T14:00:00Z"`           // 预约看房时间
	Status       int        `json:"status" example:"0"`                                    // 状态：0-待确认，1-已确认，2-已完成，3-已取消
	StatusText   string     `json:"status_text" example:"pending"`                         // 状态文本描述
	Remark       string     `json:"remark,omitempty" example:"希望周末下午看房"`                   // 备注信息
	ContactName  string     `json:"contact_name" example:"张三"`                             // 联系人姓名
	ContactPhone string     `json:"contact_phone" example:"13800138000"`                   // 联系人电话
	ConfirmTime  *time.Time `json:"confirm_time,omitempty" example:"2023-07-02T10:00:00Z"` // 确认时间
	CancelTime   *time.Time `json:"cancel_time,omitempty" example:"2023-07-02T10:00:00Z"`  // 取消时间
	CancelReason string     `json:"cancel_reason,omitempty" example:"临时有事无法看房"`            // 取消原因
	CreatedAt    time.Time  `json:"created_at" example:"2023-07-01T10:00:00Z"`             // 创建时间
}

// ViewingListResponse 预约看房列表响应DTO
type ViewingListResponse struct {
	Viewings   []ViewingResponse         `json:"viewings"`   // 预约看房列表
	Pagination common.PaginationResponse `json:"pagination"` // 分页信息
}

// GetStatusText 获取状态文本描述
func GetStatusText(status int) string {
	switch status {
	case 0:
		return "pending"
	case 1:
		return "confirmed"
	case 2:
		return "completed"
	case 3:
		return "cancelled"
	default:
		return "unknown"
	}
}
