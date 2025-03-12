package repository

import (
	"myApp/model"
)

type HouseRepository interface {
	Create(house *model.House) error
	GetByID(id uint) (*model.House, error)
	GetAll(params map[string]interface{}) ([]model.House, error)
	Update(house *model.House) error
	Delete(id uint) error
	GetHousesByLandlordID(landlordID uint) ([]model.House, error)
	IncrementViewCount(id uint) error
}

type houseRepository struct{}

func NewHouseRepository() HouseRepository {
	return &houseRepository{}
}

func (r *houseRepository) Create(house *model.House) error {
	return model.GetDB().Create(house).Error
}

func (r *houseRepository) GetByID(id uint) (*model.House, error) {
	var house model.House
	if err := model.GetDB().First(&house, id).Error; err != nil {
		return nil, err
	}
	return &house, nil
}

func (r *houseRepository) GetAll(params map[string]interface{}) ([]model.House, error) {
	var houses []model.House
	db := model.GetDB()

	// 根据参数构建查询条件
	if params != nil {
		if status, ok := params["status"].(int); ok {
			db = db.Where("status = ?", status)
		}
		if landlordID, ok := params["landlord_id"].(uint); ok {
			db = db.Where("landlord_id = ?", landlordID)
		}
		if minPrice, ok := params["min_price"].(float64); ok {
			db = db.Where("rent_price >= ?", minPrice)
		}
		if maxPrice, ok := params["max_price"].(float64); ok {
			db = db.Where("rent_price <= ?", maxPrice)
		}
		if rooms, ok := params["rooms"].(int); ok {
			db = db.Where("rooms = ?", rooms)
		}
		if houseType, ok := params["house_type"].(int); ok {
			db = db.Where("house_type = ?", houseType)
		}
		if keyword, ok := params["keyword"].(string); ok && keyword != "" {
			keyword = "%" + keyword + "%"
			db = db.Where("title LIKE ? OR description LIKE ? OR address LIKE ?", keyword, keyword, keyword)
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

	if err := db.Find(&houses).Error; err != nil {
		return nil, err
	}
	return houses, nil
}

func (r *houseRepository) Update(house *model.House) error {
	return model.GetDB().Save(house).Error
}

func (r *houseRepository) Delete(id uint) error {
	return model.GetDB().Delete(&model.House{}, id).Error
}

func (r *houseRepository) GetHousesByLandlordID(landlordID uint) ([]model.House, error) {
	var houses []model.House
	if err := model.GetDB().Where("landlord_id = ?", landlordID).Find(&houses).Error; err != nil {
		return nil, err
	}
	return houses, nil
}

func (r *houseRepository) IncrementViewCount(id uint) error {
	return model.GetDB().Model(&model.House{}).Where("id = ?", id).UpdateColumn("view_count", model.GetDB().Raw("view_count + 1")).Error
}