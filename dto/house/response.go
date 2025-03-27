package house

import (
	"time"
)

// 房源基本信息DTO
type BasicInfoDTO struct {
	ID         uint      `json:"id"`          // 房源ID
	Title      string    `json:"title"`       // 房源标题
	Address    string    `json:"address"`     // 房源地址
	Area       float64   `json:"area"`        // 房屋面积(平方米)
	Rooms      int       `json:"rooms"`       // 房间数
	Halls      int       `json:"halls"`       // 客厅数
	Bathrooms  int       `json:"bathrooms"`   // 卫生间数
	RentPrice  float64   `json:"rent_price"`  // 租金(元/月)
	HouseType  int       `json:"house_type"`  // 房屋类型
	Decoration int       `json:"decoration"`  // 装修情况
	Images     string    `json:"images"`      // 房源图片URL
	LandlordID uint      `json:"landlord_id"` // 房东ID
	Status     int       `json:"status"`      // 状态
	ViewCount  int       `json:"view_count"`  // 浏览次数
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
}

// 房源详细信息DTO
type DetailDTO struct {
	ID          uint      `json:"id"`           // 房源ID
	Title       string    `json:"title"`        // 房源标题
	Description string    `json:"description"`  // 房源描述
	Address     string    `json:"address"`      // 房源地址
	Area        float64   `json:"area"`         // 房屋面积(平方米)
	Floor       int       `json:"floor"`        // 所在楼层
	TotalFloor  int       `json:"total_floor"`  // 总楼层
	Rooms       int       `json:"rooms"`        // 房间数
	Halls       int       `json:"halls"`        // 客厅数
	Bathrooms   int       `json:"bathrooms"`    // 卫生间数
	RentPrice   float64   `json:"rent_price"`   // 租金(元/月)
	Deposit     float64   `json:"deposit"`      // 押金(元)
	PaymentType int       `json:"payment_type"` // 支付方式
	HouseType   int       `json:"house_type"`   // 房屋类型
	Orientation string    `json:"orientation"`  // 朝向
	Decoration  int       `json:"decoration"`   // 装修情况
	Facilities  string    `json:"facilities"`   // 配套设施
	Status      int       `json:"status"`       // 状态
	LandlordID  uint      `json:"landlord_id"`  // 房东ID
	Images      string    `json:"images"`       // 房源图片URL
	Latitude    float64   `json:"latitude"`     // 纬度
	Longitude   float64   `json:"longitude"`    // 经度
	IsElevator  bool      `json:"is_elevator"`  // 是否有电梯
	ViewCount   int       `json:"view_count"`   // 浏览次数
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
}

// 房源列表响应DTO
type ListResponse struct {
	Total int            `json:"total"` // 总数
	List  []BasicInfoDTO `json:"list"`  // 列表
}
