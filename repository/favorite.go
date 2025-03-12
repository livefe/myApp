package repository

import (
	"myApp/model"
)

type FavoriteRepository interface {
	Create(favorite *model.Favorite) error
	GetByID(id uint) (*model.Favorite, error)
	GetAll(params map[string]interface{}) ([]model.Favorite, error)
	Update(favorite *model.Favorite) error
	Delete(id uint) error
	GetFavoritesByUserID(userID uint) ([]model.Favorite, error)
	IsFavorite(userID, houseID uint) (bool, error)
	DeleteByUserAndHouse(userID, houseID uint) error
}

type favoriteRepository struct{}

func NewFavoriteRepository() FavoriteRepository {
	return &favoriteRepository{}
}

func (r *favoriteRepository) Create(favorite *model.Favorite) error {
	return model.GetDB().Create(favorite).Error
}

func (r *favoriteRepository) GetByID(id uint) (*model.Favorite, error) {
	var favorite model.Favorite
	if err := model.GetDB().First(&favorite, id).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *favoriteRepository) GetAll(params map[string]interface{}) ([]model.Favorite, error) {
	var favorites []model.Favorite
	db := model.GetDB()

	// 根据参数构建查询条件
	if params != nil {
		if userID, ok := params["user_id"].(uint); ok {
			db = db.Where("user_id = ?", userID)
		}
		if houseID, ok := params["house_id"].(uint); ok {
			db = db.Where("house_id = ?", houseID)
		}
	}

	// 排序
	if orderBy, ok := params["order_by"].(string); ok && orderBy != "" {
		db = db.Order(orderBy)
	} else {
		db = db.Order("created_at DESC")
	}

	// 分页
	if limit, ok := params["limit"].(int); ok && limit > 0 {
		db = db.Limit(limit)
		if offset, ok := params["offset"].(int); ok && offset >= 0 {
			db = db.Offset(offset)
		}
	}

	if err := db.Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (r *favoriteRepository) Update(favorite *model.Favorite) error {
	return model.GetDB().Save(favorite).Error
}

func (r *favoriteRepository) Delete(id uint) error {
	return model.GetDB().Delete(&model.Favorite{}, id).Error
}

func (r *favoriteRepository) GetFavoritesByUserID(userID uint) ([]model.Favorite, error) {
	var favorites []model.Favorite
	if err := model.GetDB().Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (r *favoriteRepository) IsFavorite(userID, houseID uint) (bool, error) {
	var count int64
	err := model.GetDB().Model(&model.Favorite{}).Where("user_id = ? AND house_id = ?", userID, houseID).Count(&count).Error
	return count > 0, err
}

func (r *favoriteRepository) DeleteByUserAndHouse(userID, houseID uint) error {
	return model.GetDB().Where("user_id = ? AND house_id = ?", userID, houseID).Delete(&model.Favorite{}).Error
}