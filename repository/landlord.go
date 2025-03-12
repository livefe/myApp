package repository

import (
	"myApp/model"
)

type LandlordRepository interface {
	Create(landlord *model.Landlord) error
	FindByID(id uint) (*model.Landlord, error)
	FindByUserID(userID uint) (*model.Landlord, error)
	Update(landlord *model.Landlord) error
	Delete(id uint) error
}

type landlordRepository struct{}

func NewLandlordRepository() LandlordRepository {
	return &landlordRepository{}
}

func (r *landlordRepository) Create(landlord *model.Landlord) error {
	return model.GetDB().Create(landlord).Error
}

func (r *landlordRepository) FindByID(id uint) (*model.Landlord, error) {
	var landlord model.Landlord
	if err := model.GetDB().First(&landlord, id).Error; err != nil {
		return nil, err
	}
	return &landlord, nil
}

func (r *landlordRepository) FindByUserID(userID uint) (*model.Landlord, error) {
	var landlord model.Landlord
	if err := model.GetDB().Where("user_id = ?", userID).First(&landlord).Error; err != nil {
		return nil, err
	}
	return &landlord, nil
}

func (r *landlordRepository) Update(landlord *model.Landlord) error {
	return model.GetDB().Save(landlord).Error
}

func (r *landlordRepository) Delete(id uint) error {
	return model.GetDB().Delete(&model.Landlord{}, id).Error
}