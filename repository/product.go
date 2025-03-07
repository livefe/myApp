package repository

import (
	"myApp/model"
)

type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(id uint) (*model.Product, error)
	GetAll(params map[string]interface{}) ([]model.Product, error)
	Update(product *model.Product) error
	Delete(id uint) error
	GetProductsByCreatorID(creatorID uint) ([]model.Product, error)
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (r *productRepository) Create(product *model.Product) error {
	return model.GetDB().Create(product).Error
}

func (r *productRepository) GetByID(id uint) (*model.Product, error) {
	var product model.Product
	err := model.GetDB().First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetAll(params map[string]interface{}) ([]model.Product, error) {
	var products []model.Product
	query := model.GetDB()

	if communityID, ok := params["community_id"].(uint); ok {
		query = query.Where("community_id = ?", communityID)
	}

	if categoryID, ok := params["category_id"].(uint); ok {
		query = query.Where("category_id = ?", categoryID)
	}

	if status, ok := params["status"].(int); ok {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Update(product *model.Product) error {
	return model.GetDB().Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return model.GetDB().Delete(&model.Product{}, id).Error
}

func (r *productRepository) GetProductsByCreatorID(creatorID uint) ([]model.Product, error) {
	var products []model.Product
	err := model.GetDB().Where("creator_id = ?", creatorID).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// ProductCategoryRepository 定义产品分类仓库接口
type ProductCategoryRepository interface {
	Create(category *model.ProductCategory) error
	GetByID(id uint) (*model.ProductCategory, error)
	GetAll() ([]model.ProductCategory, error)
	Update(category *model.ProductCategory) error
	Delete(id uint) error
	GetChildCategories(parentID uint) ([]model.ProductCategory, error)
}

type productCategoryRepository struct{}

func NewProductCategoryRepository() ProductCategoryRepository {
	return &productCategoryRepository{}
}

func (r *productCategoryRepository) Create(category *model.ProductCategory) error {
	return model.GetDB().Create(category).Error
}

func (r *productCategoryRepository) GetByID(id uint) (*model.ProductCategory, error) {
	var category model.ProductCategory
	err := model.GetDB().First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *productCategoryRepository) GetAll() ([]model.ProductCategory, error) {
	var categories []model.ProductCategory
	err := model.GetDB().Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *productCategoryRepository) Update(category *model.ProductCategory) error {
	return model.GetDB().Save(category).Error
}

func (r *productCategoryRepository) Delete(id uint) error {
	return model.GetDB().Delete(&model.ProductCategory{}, id).Error
}

func (r *productCategoryRepository) GetChildCategories(parentID uint) ([]model.ProductCategory, error) {
	var categories []model.ProductCategory
	err := model.GetDB().Where("parent_id = ?", parentID).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
