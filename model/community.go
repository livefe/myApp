package model

type Community struct {
	BaseModel
	Name         string `gorm:"size:100;not null" json:"name"`
	Description  string `gorm:"type:text" json:"description"`
	CreatorID    uint   `gorm:"index" json:"creator_id"`
	Status       int    `gorm:"default:1" json:"status"`
	MembersCount int    `gorm:"default:0" json:"members_count"`
	LogoURL      string `gorm:"size:255" json:"logo_url"`
}
