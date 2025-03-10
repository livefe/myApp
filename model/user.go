package model

import (
	"time"
)

type User struct {
	BaseModel
	Username  string     `gorm:"uniqueIndex;size:50" json:"username"`
	Password  string     `gorm:"size:100" json:"password,omitempty"`
	Email     string     `gorm:"uniqueIndex;size:100" json:"email"`
	Phone     string     `gorm:"size:20" json:"phone"`
	Avatar    string     `gorm:"size:255" json:"avatar"`
	LastLogin *time.Time `gorm:"default:null" json:"last_login"`
}
