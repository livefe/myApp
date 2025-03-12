package model

type Favorite struct {
	BaseModel
	UserID  uint `gorm:"type:int unsigned;comment:用户ID" json:"user_id"`  // 用户ID
	HouseID uint `gorm:"type:int unsigned;comment:房源ID" json:"house_id"` // 房源ID
	Notes   string `gorm:"type:text;comment:收藏备注" json:"notes"` // 收藏备注
}