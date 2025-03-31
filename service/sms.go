package service

import (
	"encoding/json"
	"errors"
	"myApp/config"
	"myApp/pkg/sms"
)

// SMSService 短信服务
type SMSService struct {
	Provider sms.SMSProvider
}

// NewSMSService 创建短信服务
func NewSMSService() (*SMSService, error) {
	// 创建短信服务提供商
	provider, err := sms.CreateSMSProvider()
	if err != nil {
		return nil, err
	}

	return &SMSService{
		Provider: provider,
	}, nil
}

// SendSMS 发送短信
// phoneNumbers: 接收短信的手机号码列表
// templateCode: 短信模板ID
// templateParam: 短信模板参数，可以是map或结构体，会被转换为JSON字符串
func (s *SMSService) SendSMS(phoneNumbers []string, templateCode string, templateParam interface{}) (*sms.SMSResponse, error) {
	if len(phoneNumbers) == 0 {
		return nil, errors.New("手机号码列表不能为空")
	}

	if templateCode == "" {
		return nil, errors.New("短信模板ID不能为空")
	}

	// 将模板参数转换为JSON字符串
	var templateParamStr string
	if templateParam != nil {
		paramBytes, err := json.Marshal(templateParam)
		if err != nil {
			return nil, errors.New("短信模板参数格式错误")
		}
		templateParamStr = string(paramBytes)
	} else {
		templateParamStr = "{}"
	}

	// 获取短信签名
	signName := config.Conf.SMS.Aliyun.SignName
	if signName == "" {
		return nil, errors.New("短信签名不能为空")
	}

	// 发送短信
	success, err := s.Provider.SendSMS(phoneNumbers, signName, templateCode, templateParamStr)
	if err != nil {
		return &sms.SMSResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		}, nil
	}

	return &sms.SMSResponse{
		Success: success,
		Message: "短信发送成功",
		Data:    map[string]interface{}{"provider": s.Provider.GetName()},
	}, nil
}

// QuerySMSStatus 查询短信发送状态
func (s *SMSService) QuerySMSStatus(phoneNumber, bizId string) (*sms.SMSResponse, error) {
	if phoneNumber == "" || bizId == "" {
		return nil, errors.New("手机号码和业务ID不能为空")
	}

	// 查询短信发送状态
	result, err := s.Provider.QuerySMSStatus(phoneNumber, bizId)
	if err != nil {
		return &sms.SMSResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		}, nil
	}

	return &sms.SMSResponse{
		Success: true,
		Message: "查询成功",
		Data:    result,
	}, nil
}
