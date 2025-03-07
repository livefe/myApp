package model

type CommunityMember struct {
	BaseModel
	CommunityID uint `gorm:"index" json:"community_id"`
	UserID      uint `gorm:"index" json:"user_id"`
	Role        int  `gorm:"default:0" json:"role"`
}