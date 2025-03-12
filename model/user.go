package model

import (
	"time"
)

type User struct {
	BaseModel
	Username  string     `gorm:"type:varchar(50);comment:用户名" json:"username"` // 用户名
	Password  string     `gorm:"type:varchar(100);comment:密码" json:"password,omitempty"` // 密码
	Phone     string     `gorm:"type:varchar(20);comment:手机号" json:"phone"` // 手机号
	Avatar    string     `gorm:"type:varchar(255);comment:头像URL" json:"avatar"` // 头像URL
	LastLogin *time.Time `gorm:"type:datetime;default:null;comment:最后登录时间" json:"last_login"` // 最后登录时间
	RealName  string     `gorm:"type:varchar(50);comment:真实姓名" json:"real_name"` // 真实姓名
	IdCard    string     `gorm:"type:varchar(18);comment:身份证号" json:"id_card"` // 身份证号
	Email     string     `gorm:"type:varchar(100);comment:电子邮箱" json:"email"` // 电子邮箱
	UserType  int        `gorm:"type:tinyint;default:0;comment:用户类型：0-普通用户，1-房东，2-管理员" json:"user_type"` // 用户类型：0-普通用户，1-房东，2-管理员
}
