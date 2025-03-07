package model

type Product struct {
	BaseModel
	Name        string  `gorm:"size:100;not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int     `gorm:"not null;default:0" json:"stock"`
	Status      int     `gorm:"default:1" json:"status"`
	CategoryID  uint    `gorm:"index" json:"category_id"`
	CommunityID uint    `gorm:"index" json:"community_id"`
	CreatorID   uint    `gorm:"index" json:"creator_id"`
	ImageURL    string  `gorm:"size:255" json:"image_url"`
}