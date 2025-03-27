package repository

import (
	"myApp/model"

	"gorm.io/gorm"
)

type ViewingRepository interface {
	Create(viewing *model.Viewing) error
	GetByID(id uint) (*model.Viewing, error)
	GetAll(params map[string]interface{}) ([]model.Viewing, error)
	Update(viewing *model.Viewing) error
	Delete(id uint) error
	GetViewingsByUserID(userID uint) ([]model.Viewing, error)
	GetViewingsByHouseID(houseID uint) ([]model.Viewing, error)
	UpdateStatus(id uint, status int) error
}

type viewingRepository struct {
	db *gorm.DB
}

func NewViewingRepository() ViewingRepository {
	return &viewingRepository{
		db: model.GetDB(),
	}
}

func (r *viewingRepository) Create(viewing *model.Viewing) error {
	return r.db.Create(viewing).Error
}

func (r *viewingRepository) GetByID(id uint) (*model.Viewing, error) {
	var viewing model.Viewing
	if err := r.db.First(&viewing, id).Error; err != nil {
		return nil, err
	}
	return &viewing, nil
}

func (r *viewingRepository) GetAll(params map[string]interface{}) ([]model.Viewing, error) {
	var viewings []model.Viewing
	db := r.db

	// 根据参数构建查询条件
	if params != nil {
		if status, ok := params["status"].(int); ok {
			db = db.Where("status = ?", status)
		}
		if userID, ok := params["user_id"].(uint); ok {
			db = db.Where("user_id = ?", userID)
		}
		if houseID, ok := params["house_id"].(uint); ok {
			db = db.Where("house_id = ?", houseID)
		}
		if startTime, ok := params["start_time"].(string); ok {
			db = db.Where("viewing_time >= ?", startTime)
		}
		if endTime, ok := params["end_time"].(string); ok {
			db = db.Where("viewing_time <= ?", endTime)
		}
	}

	// 排序
	if orderBy, ok := params["order_by"].(string); ok && orderBy != "" {
		db = db.Order(orderBy)
	} else {
		db = db.Order("viewing_time ASC")
	}

	// 分页
	if limit, ok := params["limit"].(int); ok && limit > 0 {
		db = db.Limit(limit)
		if offset, ok := params["offset"].(int); ok && offset >= 0 {
			db = db.Offset(offset)
		}
	}

	if err := db.Find(&viewings).Error; err != nil {
		return nil, err
	}
	return viewings, nil
}

func (r *viewingRepository) Update(viewing *model.Viewing) error {
	return r.db.Save(viewing).Error
}

func (r *viewingRepository) Delete(id uint) error {
	return r.db.Delete(&model.Viewing{}, id).Error
}

func (r *viewingRepository) GetViewingsByUserID(userID uint) ([]model.Viewing, error) {
	var viewings []model.Viewing
	if err := r.db.Where("user_id = ?", userID).Find(&viewings).Error; err != nil {
		return nil, err
	}
	return viewings, nil
}

func (r *viewingRepository) GetViewingsByHouseID(houseID uint) ([]model.Viewing, error) {
	var viewings []model.Viewing
	if err := r.db.Where("house_id = ?", houseID).Find(&viewings).Error; err != nil {
		return nil, err
	}
	return viewings, nil
}

func (r *viewingRepository) UpdateStatus(id uint, status int) error {
	return r.db.Model(&model.Viewing{}).Where("id = ?", id).Update("status", status).Error
}
