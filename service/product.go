package service

import (
	"myApp/model"
	"myApp/repository"
)

type ProductService interface {
	CreateProduct(product *model.Product) error
	GetProductByID(id uint) (*model.Product, error)
	GetAllProducts(params map[string]interface{}) ([]model.Product, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(id uint) error
	GetProductsByCreatorID(creatorID uint) ([]model.Product, error)

	// 产品分类相关方法
	CreateCategory(category *model.ProductCategory) error
	GetAllCategories() ([]model.ProductCategory, error)
	GetCategoryByID(id uint) (*model.ProductCategory, error)
	UpdateCategory(category *model.ProductCategory) error
	DeleteCategory(id uint) error
}

type ProductServiceImpl struct {
	ProductRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		ProductRepo: productRepo,
	}
}

func (s *ProductServiceImpl) CreateProduct(product *model.Product) error {
	return s.ProductRepo.Create(product)
}

func (s *ProductServiceImpl) GetProductByID(id uint) (*model.Product, error) {
	return s.ProductRepo.GetByID(id)
}

func (s *ProductServiceImpl) GetAllProducts(params map[string]interface{}) ([]model.Product, error) {
	return s.ProductRepo.GetAll(params)
}

func (s *ProductServiceImpl) UpdateProduct(product *model.Product) error {
	return s.ProductRepo.Update(product)
}

func (s *ProductServiceImpl) DeleteProduct(id uint) error {
	return s.ProductRepo.Delete(id)
}

func (s *ProductServiceImpl) GetProductsByCreatorID(creatorID uint) ([]model.Product, error) {
	return s.ProductRepo.GetProductsByCreatorID(creatorID)
}

// 产品分类相关方法实现
func (s *ProductServiceImpl) CreateCategory(category *model.ProductCategory) error {
	return s.ProductRepo.CreateCategory(category)
}

func (s *ProductServiceImpl) GetAllCategories() ([]model.ProductCategory, error) {
	return s.ProductRepo.GetAllCategories()
}

func (s *ProductServiceImpl) GetCategoryByID(id uint) (*model.ProductCategory, error) {
	return s.ProductRepo.GetCategoryByID(id)
}

func (s *ProductServiceImpl) UpdateCategory(category *model.ProductCategory) error {
	return s.ProductRepo.UpdateCategory(category)
}

func (s *ProductServiceImpl) DeleteCategory(id uint) error {
	// 应该调用专门用于删除分类的方法
	return s.ProductRepo.DeleteCategory(id)
}
