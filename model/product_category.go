package model

type ProductCategory struct {
	BaseModel
	Name        string `gorm:"size:50;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	ParentID    *uint  `gorm:"" json:"parent_id"`
}