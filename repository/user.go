package repository

import (
	"myApp/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindByPhone(phone string) ([]*model.User, error)
	Update(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: model.GetDB(),
	}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) FindByPhone(phone string) ([]*model.User, error) {
	var users []*model.User
	if err := r.db.Where("phone = ?", phone).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
