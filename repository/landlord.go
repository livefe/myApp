package repository

import (
	"myApp/model"

	"gorm.io/gorm"
)

type LandlordRepository interface {
	Create(landlord *model.Landlord) error
	FindByID(id uint) (*model.Landlord, error)
	FindByUserID(userID uint) (*model.Landlord, error)
	Update(landlord *model.Landlord) error
	Delete(id uint) error
}

type landlordRepository struct {
	db *gorm.DB
}

func NewLandlordRepository() LandlordRepository {
	return &landlordRepository{
		db: model.GetDB(),
	}
}

func (r *landlordRepository) Create(landlord *model.Landlord) error {
	return r.db.Create(landlord).Error
}

func (r *landlordRepository) FindByID(id uint) (*model.Landlord, error) {
	var landlord model.Landlord
	if err := r.db.First(&landlord, id).Error; err != nil {
		return nil, err
	}
	return &landlord, nil
}

func (r *landlordRepository) FindByUserID(userID uint) (*model.Landlord, error) {
	var landlord model.Landlord
	if err := r.db.Where("user_id = ?", userID).First(&landlord).Error; err != nil {
		return nil, err
	}
	return &landlord, nil
}

func (r *landlordRepository) Update(landlord *model.Landlord) error {
	return r.db.Save(landlord).Error
}

func (r *landlordRepository) Delete(id uint) error {
	return r.db.Delete(&model.Landlord{}, id).Error
}
