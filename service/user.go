package service

import (
	"errors"
	"myApp/model"
	"myApp/pkg/logger"
	"myApp/repository"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *model.User) (*model.User, error)
	Login(username, password string) (*model.User, error)
	GetUserProfile(id uint) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(user *model.User) (*model.User, error) {
	// 记录用户注册开始
	logger.Info("用户注册开始",
		zap.String("username", user.Username),
		zap.String("email", user.Email),
		zap.String("phone", user.Phone),
	)

	// 检查用户名是否已存在
	existingUser, _ := s.repo.FindByUsername(user.Username)
	if existingUser != nil {
		logger.Warn("用户注册失败：用户名已存在", zap.String("username", user.Username))
		return nil, errors.New("用户名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("用户注册失败：密码加密失败", zap.Error(err))
		return nil, errors.New("密码加密失败")
	}

	user.Password = string(hashedPassword)
	err = s.repo.Create(user)
	if err != nil {
		logger.Error("用户注册失败：创建用户失败", zap.Error(err))
		return nil, errors.New("创建用户失败")
	}

	// 记录用户注册成功
	logger.Info("用户注册成功",
		zap.String("username", user.Username),
		zap.Uint("user_id", user.ID),
	)

	return user, nil
}

func (s *userService) Login(username, password string) (*model.User, error) {
	// 记录登录尝试
	logger.Info("用户登录尝试", zap.String("username", username))

	user, err := s.repo.FindByUsername(username)
	if err != nil {
		logger.Warn("用户登录失败：用户不存在", zap.String("username", username))
		return nil, errors.New("用户不存在")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Warn("用户登录失败：密码错误",
			zap.String("username", username),
			zap.Uint("user_id", user.ID),
		)
		return nil, errors.New("密码错误")
	}

	// 记录登录成功
	logger.Info("用户登录成功",
		zap.String("username", username),
		zap.Uint("user_id", user.ID),
	)

	return user, nil
}

func (s *userService) GetUserProfile(id uint) (*model.User, error) {
	logger.Debug("获取用户资料", zap.Uint("user_id", id))

	user, err := s.repo.FindByID(id)
	if err != nil {
		logger.Warn("获取用户资料失败：用户不存在", zap.Uint("user_id", id))
		return nil, errors.New("用户不存在")
	}

	logger.Debug("获取用户资料成功", zap.Uint("user_id", id))
	return user, nil
}
