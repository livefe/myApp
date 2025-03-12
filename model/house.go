package model

type House struct {
	BaseModel
	Title       string  `gorm:"type:varchar(100);not null;comment:房源标题" json:"title"`       // 房源标题
	Description string  `gorm:"type:text;comment:房源描述" json:"description"`         // 房源描述
	Address     string  `gorm:"type:varchar(255);not null;comment:房源地址" json:"address"`     // 房源地址
	Area        float64 `gorm:"type:decimal(10,2);not null;comment:房屋面积(平方米)" json:"area"` // 房屋面积(平方米)
	Floor       int     `gorm:"type:int;comment:所在楼层" json:"floor"`                        // 所在楼层
	TotalFloor  int     `gorm:"type:int;comment:总楼层" json:"total_floor"`                  // 总楼层
	Rooms       int     `gorm:"type:int;not null;comment:房间数" json:"rooms"`                // 房间数
	Halls       int     `gorm:"type:int;not null;comment:客厅数" json:"halls"`                // 客厅数
	Bathrooms   int     `gorm:"type:int;not null;comment:卫生间数" json:"bathrooms"`            // 卫生间数
	RentPrice   float64 `gorm:"type:decimal(10,2);not null;comment:租金(元/月)" json:"rent_price"` // 租金(元/月)
	Deposit     float64 `gorm:"type:decimal(10,2);comment:押金(元)" json:"deposit"`     // 押金(元)
	PaymentType int     `gorm:"type:tinyint;default:1;comment:支付方式：1-月付，2-季付，3-半年付，4-年付" json:"payment_type"`        // 支付方式：1-月付，2-季付，3-半年付，4-年付
	HouseType   int     `gorm:"type:tinyint;not null;comment:房屋类型：1-普通住宅，2-公寓，3-别墅，4-商铺" json:"house_type"`           // 房屋类型：1-普通住宅，2-公寓，3-别墅，4-商铺
	Orientation string  `gorm:"type:varchar(20);comment:朝向" json:"orientation"`           // 朝向
	Decoration  int     `gorm:"type:tinyint;default:1;comment:装修情况：1-简装，2-精装，3-豪装" json:"decoration"`          // 装修情况：1-简装，2-精装，3-豪装
	Facilities  string  `gorm:"type:text;comment:配套设施，JSON格式字符串" json:"facilities"`          // 配套设施，JSON格式字符串
	Status      int     `gorm:"type:tinyint;default:1;comment:状态：0-下架，1-上架" json:"status"`              // 状态：0-下架，1-上架
	LandlordID  uint    `gorm:"type:int unsigned;comment:房东ID" json:"landlord_id"`                  // 房东ID
	Images      string  `gorm:"type:text;comment:房源图片URL，JSON格式字符串" json:"images"`              // 房源图片URL，JSON格式字符串
	Latitude    float64 `gorm:"type:decimal(10,6);comment:纬度" json:"latitude"`    // 纬度
	Longitude   float64 `gorm:"type:decimal(10,6);comment:经度" json:"longitude"`   // 经度
	IsElevator  bool    `gorm:"type:tinyint(1);default:false;comment:是否有电梯" json:"is_elevator"`      // 是否有电梯
	ViewCount   int     `gorm:"type:int;default:0;comment:浏览次数" json:"view_count"`          // 浏览次数
}