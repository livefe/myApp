package service

import (
	"errors"
	"myApp/model"
	"myApp/repository"

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
	// 检查用户名是否已存在
	existingUser, _ := s.repo.FindByUsername(user.Username)
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user.Password = string(hashedPassword)
	err = s.repo.Create(user)
	if err != nil {
		return nil, errors.New("创建用户失败")
	}

	return user, nil
}

func (s *userService) Login(username, password string) (*model.User, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

func (s *userService) GetUserProfile(id uint) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}
