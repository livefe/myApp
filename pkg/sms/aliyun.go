package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"myApp/config"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
)

// AliyunSMSConfig 阿里云短信配置
type AliyunSMSConfig struct {
	AccessKeyID     string `json:"accessKeyId"`     // 阿里云AccessKey ID
	AccessKeySecret string `json:"accessKeySecret"` // 阿里云AccessKey Secret
	RegionID        string `json:"regionId"`        // 地域ID
}

// GetConfig 获取配置信息
func (c *AliyunSMSConfig) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"accessKeyId":     c.AccessKeyID,
		"accessKeySecret": c.AccessKeySecret,
		"regionId":        c.RegionID,
	}
}

// AliyunSMSProvider 阿里云短信服务提供商
type AliyunSMSProvider struct {
	Config *AliyunSMSConfig
	Client *dysmsapi.Client
}

// NewAliyunSMSProvider 创建阿里云短信服务提供商实例
func NewAliyunSMSProvider(config *AliyunSMSConfig) (*AliyunSMSProvider, error) {
	if config == nil {
		return nil, errors.New("阿里云短信配置不能为空")
	}

	// 创建阿里云短信客户端
	client, err := createAliyunClient(config)
	if err != nil {
		return nil, err
	}

	return &AliyunSMSProvider{
		Config: config,
		Client: client,
	}, nil
}

// 创建阿里云客户端
func createAliyunClient(config *AliyunSMSConfig) (*dysmsapi.Client, error) {
	// 创建配置
	clientConfig := &openapi.Config{
		AccessKeyId:     tea.String(config.AccessKeyID),
		AccessKeySecret: tea.String(config.AccessKeySecret),
	}
	// 设置Endpoint
	clientConfig.Endpoint = tea.String("dysmsapi.aliyuncs.com")

	// 创建客户端
	client, err := dysmsapi.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云短信客户端失败: %v", err)
	}

	return client, nil
}

// SendSMS 发送短信
func (p *AliyunSMSProvider) SendSMS(phoneNumbers []string, signName, templateCode, templateParam string) (bool, string, string, error) {
	if len(phoneNumbers) == 0 {
		return false, "", "", errors.New("手机号码列表不能为空")
	}

	// 将手机号码列表转换为逗号分隔的字符串
	phoneNumbersStr := ""
	for i, phone := range phoneNumbers {
		if i > 0 {
			phoneNumbersStr += ","
		}
		phoneNumbersStr += phone
	}

	// 创建发送短信请求
	sendSmsRequest := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phoneNumbersStr),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(templateParam),
	}

	// 发送短信
	response, err := p.Client.SendSms(sendSmsRequest)
	if err != nil {
		return false, "", "", fmt.Errorf("发送短信失败: %v", err)
	}

	// 检查发送结果
	if *response.Body.Code != "OK" {
		return false, "", "", fmt.Errorf("发送短信失败: %s, %s", *response.Body.Code, *response.Body.Message)
	}

	// 获取BizId和RequestId
	bizId := ""
	requestId := ""
	if response.Body.BizId != nil {
		bizId = *response.Body.BizId
	}
	if response.Body.RequestId != nil {
		requestId = *response.Body.RequestId
	}

	return true, bizId, requestId, nil
}

// QuerySMSStatus 查询短信发送状态
func (p *AliyunSMSProvider) QuerySMSStatus(phoneNumber, bizId string) (map[string]interface{}, error) {
	if phoneNumber == "" || bizId == "" {
		return nil, errors.New("手机号码和业务ID不能为空")
	}

	// 创建查询短信发送状态请求
	queryRequest := &dysmsapi.QuerySendDetailsRequest{
		PhoneNumber: tea.String(phoneNumber),
		BizId:       tea.String(bizId),
	}

	// 查询短信发送状态
	response, err := p.Client.QuerySendDetails(queryRequest)
	if err != nil {
		return nil, fmt.Errorf("查询短信发送状态失败: %v", err)
	}

	// 检查查询结果
	if *response.Body.Code != "OK" {
		return nil, fmt.Errorf("查询短信发送状态失败: %s, %s", *response.Body.Code, *response.Body.Message)
	}

	// 将查询结果转换为map
	result := make(map[string]interface{})
	result["Code"] = *response.Body.Code
	result["Message"] = *response.Body.Message
	result["TotalCount"] = *response.Body.TotalCount

	// 将SmsSendDetailDTOs转换为JSON字符串，再解析为map
	detailsBytes, err := json.Marshal(response.Body.SmsSendDetailDTOs)
	if err != nil {
		return nil, fmt.Errorf("解析短信发送详情失败: %v", err)
	}

	var details map[string]interface{}
	if err := json.Unmarshal(detailsBytes, &details); err != nil {
		return nil, fmt.Errorf("解析短信发送详情失败: %v", err)
	}

	result["Details"] = details

	return result, nil
}

// GetName 获取短信服务提供商名称
func (p *AliyunSMSProvider) GetName() string {
	return "Aliyun"
}

// GetSignName 获取短信签名
func (p *AliyunSMSProvider) GetSignName() string {
	return config.Conf.SMS.Aliyun.SignName
}

// GetTemplateCode 获取短信模板ID
func (p *AliyunSMSProvider) GetTemplateCode() string {
	return config.Conf.SMS.Aliyun.TemplateCode
}
