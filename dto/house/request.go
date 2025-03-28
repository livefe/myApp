package house

import (
	"myApp/dto/common"

	"github.com/go-playground/validator/v10"
)

// 创建房源请求DTO
type CreateRequest struct {
	Title       string  `json:"title" binding:"required" example:"精装修两居室"`                                // 房源标题
	Description string  `json:"description" binding:"required" example:"位于市中心的精装修两居室，交通便利"`               // 房源描述
	Address     string  `json:"address" binding:"required" example:"北京市朝阳区建国路1号"`                         // 房源地址
	Area        float64 `json:"area" binding:"required,gt=0" example:"80.5"`                              // 房屋面积(平方米)
	Floor       int     `json:"floor" binding:"required,gte=0" example:"8"`                               // 所在楼层
	TotalFloor  int     `json:"total_floor" binding:"required,gt=0" example:"20"`                         // 总楼层
	Rooms       int     `json:"rooms" binding:"required,gte=1" example:"2"`                               // 房间数
	Halls       int     `json:"halls" binding:"required,gte=0" example:"1"`                               // 客厅数
	Bathrooms   int     `json:"bathrooms" binding:"required,gte=1" example:"1"`                           // 卫生间数
	RentPrice   float64 `json:"rent_price" binding:"required,gt=0" example:"5000"`                        // 租金(元/月)
	Deposit     float64 `json:"deposit" binding:"omitempty,gte=0" example:"10000"`                        // 押金(元)
	PaymentType int     `json:"payment_type" binding:"required,oneof=1 2 3 4" example:"1"`                // 支付方式：1-月付，2-季付，3-半年付，4-年付
	HouseType   int     `json:"house_type" binding:"required,oneof=1 2 3 4" example:"1"`                  // 房屋类型：1-普通住宅，2-公寓，3-别墅，4-商铺
	Orientation string  `json:"orientation" binding:"omitempty" example:"南"`                              // 朝向
	Decoration  int     `json:"decoration" binding:"required,oneof=1 2 3" example:"2"`                    // 装修情况：1-简装，2-精装，3-豪装
	Facilities  string  `json:"facilities" binding:"omitempty" example:"[\"空调\",\"热水器\",\"冰箱\",\"洗衣机\"]"` // 配套设施，JSON格式字符串
	Images      string  `json:"images" binding:"omitempty" example:"[\"http://example.com/img1.jpg\"]"`   // 房源图片URL，JSON格式字符串
	Latitude    float64 `json:"latitude" binding:"omitempty" example:"39.9087243"`                        // 纬度
	Longitude   float64 `json:"longitude" binding:"omitempty" example:"116.3952859"`                      // 经度
	IsElevator  bool    `json:"is_elevator" example:"true"`                                               // 是否有电梯
}

// 更新房源请求DTO
type UpdateRequest struct {
	Title       string  `json:"title" binding:"omitempty" example:"精装修两居室"`                               // 房源标题
	Description string  `json:"description" binding:"omitempty" example:"位于市中心的精装修两居室，交通便利"`              // 房源描述
	Address     string  `json:"address" binding:"omitempty" example:"北京市朝阳区建国路1号"`                        // 房源地址
	Area        float64 `json:"area" binding:"omitempty,gt=0" example:"80.5"`                             // 房屋面积(平方米)
	Floor       int     `json:"floor" binding:"omitempty,gte=0" example:"8"`                              // 所在楼层
	TotalFloor  int     `json:"total_floor" binding:"omitempty,gt=0" example:"20"`                        // 总楼层
	Rooms       int     `json:"rooms" binding:"omitempty,gte=1" example:"2"`                              // 房间数
	Halls       int     `json:"halls" binding:"omitempty,gte=0" example:"1"`                              // 客厅数
	Bathrooms   int     `json:"bathrooms" binding:"omitempty,gte=1" example:"1"`                          // 卫生间数
	RentPrice   float64 `json:"rent_price" binding:"omitempty,gt=0" example:"5000"`                       // 租金(元/月)
	Deposit     float64 `json:"deposit" binding:"omitempty,gte=0" example:"10000"`                        // 押金(元)
	PaymentType int     `json:"payment_type" binding:"omitempty,oneof=1 2 3 4" example:"1"`               // 支付方式：1-月付，2-季付，3-半年付，4-年付
	HouseType   int     `json:"house_type" binding:"omitempty,oneof=1 2 3 4" example:"1"`                 // 房屋类型：1-普通住宅，2-公寓，3-别墅，4-商铺
	Orientation string  `json:"orientation" binding:"omitempty" example:"南"`                              // 朝向
	Decoration  int     `json:"decoration" binding:"omitempty,oneof=1 2 3" example:"2"`                   // 装修情况：1-简装，2-精装，3-豪装
	Facilities  string  `json:"facilities" binding:"omitempty" example:"[\"空调\",\"热水器\",\"冰箱\",\"洗衣机\"]"` // 配套设施，JSON格式字符串
	Images      string  `json:"images" binding:"omitempty" example:"[\"http://example.com/img1.jpg\"]"`   // 房源图片URL，JSON格式字符串
	Latitude    float64 `json:"latitude" binding:"omitempty" example:"39.9087243"`                        // 纬度
	Longitude   float64 `json:"longitude" binding:"omitempty" example:"116.3952859"`                      // 经度
	IsElevator  bool    `json:"is_elevator" example:"true"`                                               // 是否有电梯
	Status      int     `json:"status" binding:"omitempty,oneof=0 1" example:"1"`                         // 状态：0-下架，1-上架
}

// 房源查询请求DTO
type QueryRequest struct {
	Keyword                      string  `json:"keyword" form:"keyword" example:"精装修"`          // 关键词
	Status                       int     `json:"status" form:"status" example:"1"`              // 状态：0-下架，1-上架
	LandlordID                   uint    `json:"landlord_id" form:"landlord_id" example:"1"`    // 房东ID
	MinPrice                     float64 `json:"min_price" form:"min_price" example:"3000"`     // 最低价格
	MaxPrice                     float64 `json:"max_price" form:"max_price" example:"6000"`     // 最高价格
	MinArea                      float64 `json:"min_area" form:"min_area" example:"60"`         // 最小面积
	MaxArea                      float64 `json:"max_area" form:"max_area" example:"100"`        // 最大面积
	Rooms                        int     `json:"rooms" form:"rooms" example:"2"`                // 房间数
	HouseType                    int     `json:"house_type" form:"house_type" example:"1"`      // 房屋类型
	Decoration                   int     `json:"decoration" form:"decoration" example:"2"`      // 装修情况
	IsElevator                   bool    `json:"is_elevator" form:"is_elevator" example:"true"` // 是否有电梯
	Orientation                  string  `json:"orientation" form:"orientation" example:"南"`    // 朝向
	PaymentType                  int     `json:"payment_type" form:"payment_type" example:"1"`  // 支付方式
	common.PaginationSortRequest         // 分页和排序参数
}

// ValidateCreateRequest 验证创建房源请求
func ValidateCreateRequest(req CreateRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}

// ValidateUpdateRequest 验证更新房源请求
func ValidateUpdateRequest(req UpdateRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}
