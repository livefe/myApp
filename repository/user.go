package repository

import (
	"myApp/model"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	Update(user *model.User) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *model.User) error {
	return model.GetDB().Create(user).Error
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := model.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := model.GetDB().First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return model.GetDB().Save(user).Error
}
