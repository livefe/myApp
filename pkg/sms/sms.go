package sms

// SMSProvider 定义短信服务提供商的通用接口
type SMSProvider interface {
	// SendSMS 发送短信的通用方法
	// phoneNumbers: 接收短信的手机号码列表
	// signName: 短信签名
	// templateCode: 短信模板ID
	// templateParam: 短信模板参数，JSON格式字符串
	// 返回发送结果和错误信息
	SendSMS(phoneNumbers []string, signName, templateCode, templateParam string) (bool, error)

	// QuerySMSStatus 查询短信发送状态
	// phoneNumber: 手机号码
	// bizId: 发送回执ID，可以根据发送回执ID查询具体的发送状态
	// 返回查询结果和错误信息
	QuerySMSStatus(phoneNumber, bizId string) (map[string]interface{}, error)

	// GetName 获取短信服务提供商名称
	GetName() string
}

// SMSConfig 短信配置接口
type SMSConfig interface {
	// GetConfig 获取配置信息
	GetConfig() map[string]interface{}
}

// SMSResponse 短信发送响应结构
type SMSResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
