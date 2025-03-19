package model

import (
	"time"
)

type Viewing struct {
	BaseModel
	HouseID     uint       `gorm:"type:int unsigned;comment:房源ID" json:"house_id"`                // 房源ID
	UserID      uint       `gorm:"type:int unsigned;comment:用户ID" json:"user_id"`                 // 用户ID
	ViewingTime time.Time  `gorm:"type:datetime;not null;comment:预约看房时间" json:"viewing_time"` // 预约看房时间
	Status      int        `gorm:"type:tinyint;default:0;comment:状态：0-待确认，1-已确认，2-已完成，3-已取消" json:"status"`          // 状态：0-待确认，1-已确认，2-已完成，3-已取消
	Remark      string     `gorm:"type:text;comment:备注信息" json:"remark"`          // 备注信息
	ContactName string     `gorm:"type:varchar(50);comment:联系人姓名" json:"contact_name"`      // 联系人姓名
	ContactPhone string    `gorm:"type:varchar(20);comment:联系人电话" json:"contact_phone"`     // 联系人电话
	ConfirmTime *time.Time `gorm:"type:datetime;default:null;comment:确认时间" json:"confirm_time"` // 确认时间
	CancelTime  *time.Time `gorm:"type:datetime;default:null;comment:取消时间" json:"cancel_time"`  // 取消时间
	CancelReason string    `gorm:"type:text;comment:取消原因" json:"cancel_reason"`   // 取消原因
}

// 预约看房状态常量
const (
	ViewingPending   = 0 // 待确认
	ViewingConfirmed = 1 // 已确认
	ViewingCompleted = 2 // 已完成
	ViewingCancelled = 3 // 已取消
)