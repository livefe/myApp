package viewing

import (
	"time"
	"myApp/dto/common"

	"github.com/go-playground/validator/v10"
)

// 创建预约看房请求DTO
type CreateRequest struct {
	HouseID      uint      `json:"house_id" binding:"required" example:"1"`                                      // 房源ID
	ViewDate     time.Time `json:"view_date" binding:"required,gtfield=time.Now" example:"2023-07-01T14:00:00Z"` // 预约看房时间
	Message      string    `json:"message" binding:"omitempty" example:"希望周末下午看房，最好能详细介绍下周边设施"`                  // 备注信息
	ContactName  string    `json:"contact_name" binding:"required" example:"张三"`                                 // 联系人姓名
	ContactPhone string    `json:"contact_phone" binding:"required,len=11" example:"13800138000"`                // 联系人电话
}

// 更新预约看房状态请求DTO
type UpdateStatusRequest struct {
	Status       int    `json:"status" binding:"required,oneof=0 1 2 3" example:"1"`  // 状态：0-待确认，1-已确认，2-已完成，3-已取消
	CancelReason string `json:"cancel_reason" binding:"omitempty" example:"临时有事无法看房"` // 取消原因（仅当状态为已取消时需要）
}

// 预约看房查询请求DTO
type QueryRequest struct {
	UserID                       uint      `json:"user_id" form:"user_id" example:"1"`                          // 用户ID
	HouseID                      uint      `json:"house_id" form:"house_id" example:"1"`                        // 房源ID
	Status                       int       `json:"status" form:"status" example:"0"`                            // 状态
	StartDate                    time.Time `json:"start_date" form:"start_date" example:"2023-07-01T00:00:00Z"` // 开始日期
	EndDate                      time.Time `json:"end_date" form:"end_date" example:"2023-07-31T23:59:59Z"`     // 结束日期
	common.PaginationSortRequest           // 分页和排序参数
}

// ValidateCreateRequest 验证创建预约看房请求
func ValidateCreateRequest(req CreateRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}

// ValidateUpdateStatusRequest 验证更新预约看房状态请求
func ValidateUpdateStatusRequest(req UpdateStatusRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}
