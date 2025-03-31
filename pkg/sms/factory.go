package sms

import (
	"errors"
	"fmt"
	"myApp/config"
)

// SMSFactory 短信服务工厂
type SMSFactory struct {
	Providers map[string]SMSProvider
}

// NewSMSFactory 创建短信服务工厂
func NewSMSFactory() *SMSFactory {
	return &SMSFactory{
		Providers: make(map[string]SMSProvider),
	}
}

// RegisterProvider 注册短信服务提供商
func (f *SMSFactory) RegisterProvider(name string, provider SMSProvider) {
	f.Providers[name] = provider
}

// GetProvider 获取短信服务提供商
func (f *SMSFactory) GetProvider(name string) (SMSProvider, error) {
	provider, ok := f.Providers[name]
	if !ok {
		return nil, fmt.Errorf("短信服务提供商 %s 未注册", name)
	}
	return provider, nil
}

// CreateSMSProvider 根据配置创建短信服务提供商
func CreateSMSProvider() (SMSProvider, error) {
	// 获取配置
	smsConfig := config.Conf.SMS

	// 根据配置的提供商类型创建对应的短信服务提供商
	switch smsConfig.Provider {
	case "aliyun":
		// 创建阿里云短信服务提供商
		aliyunConfig := &AliyunSMSConfig{
			AccessKeyID:     smsConfig.Aliyun.AccessKeyID,
			AccessKeySecret: smsConfig.Aliyun.AccessKeySecret,
			RegionID:        smsConfig.Aliyun.RegionID,
		}
		return NewAliyunSMSProvider(aliyunConfig)
	default:
		return nil, errors.New("不支持的短信服务提供商类型")
	}
}
