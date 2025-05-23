package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"myApp/model"
	"myApp/pkg/redis"
	"myApp/pkg/sms"
	"myApp/repository"
	"strconv"
	"time"
)

// 定义常量
const (
	SMSCodePrefix = "sms:code:" // Redis中存储短信验证码的前缀
	SMSCodeExpire = 300         // 短信验证码有效期（秒）
	SMSCodeLength = 6           // 短信验证码长度
)

// SMSCodeService 短信验证码服务接口
type SMSCodeService interface {
	SendCode(phone string, ipAddress, userAgent string) (bool, error) // 发送验证码
	VerifyCode(phone, code string) (bool, error)                      // 验证验证码
	LoginByCode(phone, code string) (*model.User, error)              // 通过验证码登录
}

// smsCodeService 短信验证码服务实现
type smsCodeService struct {
	userRepo      repository.UserRepository
	smsRecordRepo repository.SMSRecordRepository
}

// NewSMSCodeService 创建短信验证码服务实例
func NewSMSCodeService(userRepo repository.UserRepository, smsRecordRepo repository.SMSRecordRepository) SMSCodeService {
	return &smsCodeService{
		userRepo:      userRepo,
		smsRecordRepo: smsRecordRepo,
	}
}

// generateCode 生成随机验证码
func (s *smsCodeService) generateCode() string {
	// 使用crypto/rand包生成更安全的随机数
	code := ""
	for i := 0; i < SMSCodeLength; i++ {
		// 生成0-9之间的随机数
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			// 如果生成失败，回退到使用math/rand
			code += strconv.Itoa(int(time.Now().UnixNano() % 10))
			continue
		}
		code += strconv.FormatInt(num.Int64(), 10)
	}

	return code
}

// SendCode 发送短信验证码
func (s *smsCodeService) SendCode(phone string, ipAddress, userAgent string) (bool, error) {
	if phone == "" {
		return false, errors.New("手机号不能为空")
	}

	// 先检查Redis中是否存在未过期的验证码
	key := SMSCodePrefix + phone
	existingCode, err := redis.Get(key)
	// 如果验证码存在且未过期，直接返回成功
	if err == nil && existingCode != "" {
		return true, nil
	}

	// 生成验证码
	code := s.generateCode()

	// 存储验证码到Redis
	err = redis.Set(key, code, time.Duration(SMSCodeExpire)*time.Second)
	if err != nil {
		return false, fmt.Errorf("存储验证码失败: %v", err)
	}

	// 创建短信服务提供商
	provider, err := sms.CreateSMSProvider()
	if err != nil {
		return false, fmt.Errorf("创建短信服务提供商失败: %v", err)
	}

	// 构建短信模板参数
	templateParam := fmt.Sprintf(`{"code":"%s"}`, code)

	// 发送短信
	success, bizId, requestId, err := provider.SendSMS(
		[]string{phone},
		provider.GetSignName(),
		provider.GetTemplateCode(),
		templateParam,
	)

	// 创建短信记录
	smsRecord := &model.SMSRecord{
		Phone:      phone,
		Code:       code,
		TemplateID: provider.GetTemplateCode(),
		Content:    templateParam,
		Status:     success,
		Provider:   provider.GetName(),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		BizId:      bizId,
		RequestId:  requestId,
	}

	// 如果发送失败，记录失败原因
	if err != nil {
		smsRecord.Status = false
		smsRecord.FailReason = err.Error()
		// 记录短信发送失败日志，但不影响主流程返回
		_ = s.smsRecordRepo.Create(smsRecord)
		return false, fmt.Errorf("发送短信失败: %v", err)
	}

	// 记录短信发送成功日志
	_ = s.smsRecordRepo.Create(smsRecord)

	return success, nil
}

// VerifyCode 验证短信验证码
func (s *smsCodeService) VerifyCode(phone, code string) (bool, error) {
	if phone == "" || code == "" {
		return false, errors.New("手机号和验证码不能为空")
	}

	// 从Redis获取存储的验证码
	key := SMSCodePrefix + phone
	storedCode, err := redis.Get(key)

	// 如果获取失败或验证码不存在
	if err != nil {
		if err == redis.Nil {
			return false, errors.New("验证码已过期或不存在")
		}
		return false, fmt.Errorf("获取验证码失败: %v", err)
	}

	// 验证码比对
	if storedCode != code {
		return false, errors.New("验证码错误")
	}

	// 验证成功后删除验证码，防止重复使用
	_ = redis.Delete(key)

	return true, nil
}

// LoginByCode 通过验证码登录
func (s *smsCodeService) LoginByCode(phone, code string) (*model.User, error) {
	// 验证验证码
	valid, err := s.VerifyCode(phone, code)
	if !valid || err != nil {
		return nil, errors.New("验证码验证失败")
	}

	// 查找用户
	user, err := s.findOrCreateUserByPhone(phone)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("更新用户登录时间失败: %v", err)
	}

	return user, nil
}

// findOrCreateUserByPhone 根据手机号查找用户，如果不存在则创建
func (s *smsCodeService) findOrCreateUserByPhone(phone string) (*model.User, error) {
	// 尝试通过手机号查找用户
	users, err := s.userRepo.FindByPhone(phone)
	if err == nil && len(users) > 0 {
		// 用户存在，返回第一个匹配的用户
		return users[0], nil
	}

	// 用户不存在，创建新用户
	newUser := &model.User{
		Phone:    phone,
		Username: phone, // 默认使用手机号作为用户名
		UserType: 0,     // 默认为普通用户
	}

	// 创建用户
	err = s.userRepo.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	return newUser, nil
}
