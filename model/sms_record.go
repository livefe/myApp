package model

// SMSRecord 短信记录模型
// 用于记录短信验证码的发送记录，包括接收手机号、验证码内容、发送时间、发送状态等信息
type SMSRecord struct {
	BaseModel
	Phone      string `gorm:"type:varchar(20);index;comment:手机号码" json:"phone"`
	Code       string `gorm:"type:varchar(10);comment:验证码内容" json:"code"`
	TemplateID string `gorm:"type:varchar(50);comment:短信模板ID" json:"template_id"`
	Content    string `gorm:"type:varchar(255);comment:短信内容" json:"content"`
	Status     bool   `gorm:"type:tinyint(1);comment:发送状态(0失败,1成功)" json:"status"`
	FailReason string `gorm:"type:varchar(255);comment:失败原因" json:"fail_reason"`
	Provider   string `gorm:"type:varchar(50);comment:短信服务提供商" json:"provider"`
	IPAddress  string `gorm:"type:varchar(50);comment:请求IP地址" json:"ip_address"`
	UserAgent  string `gorm:"type:varchar(255);comment:用户代理" json:"user_agent"`
	BizId      string `gorm:"type:varchar(50);comment:发送回执ID" json:"biz_id"`
	RequestId  string `gorm:"type:varchar(50);comment:请求ID" json:"request_id"`
}

// TableName 指定表名
func (SMSRecord) TableName() string {
	return "sms_records"
}
