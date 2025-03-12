package model

type Landlord struct {
	BaseModel
	UserID       uint   `gorm:"type:int unsigned;comment:关联的用户ID" json:"user_id"`              // 关联的用户ID
	Verified     bool   `gorm:"type:tinyint(1);default:false;comment:是否已认证" json:"verified"`           // 是否已认证
	IdCardFront  string `gorm:"type:varchar(255);comment:身份证正面照片URL" json:"id_card_front"`          // 身份证正面照片URL
	IdCardBack   string `gorm:"type:varchar(255);comment:身份证背面照片URL" json:"id_card_back"`           // 身份证背面照片URL
	BankAccount  string `gorm:"type:varchar(50);comment:银行账号" json:"bank_account"`            // 银行账号
	BankName     string `gorm:"type:varchar(100);comment:开户行名称" json:"bank_name"`              // 开户行名称
	AccountName  string `gorm:"type:varchar(50);comment:开户人姓名" json:"account_name"`            // 开户人姓名
	Introduction string `gorm:"type:text;comment:房东自我介绍" json:"introduction"`           // 房东自我介绍
	Rating       float64 `gorm:"type:decimal(2,1);default:5.0;comment:房东评分" json:"rating"` // 房东评分
}