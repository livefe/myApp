package service

import (
	"myApp/model"
	"myApp/repository"
)

type ProductService struct {
	ProductRepo         repository.ProductRepository
	ProductCategoryRepo repository.ProductCategoryRepository
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.ProductCategoryRepository) *ProductService {
	return &ProductService{
		ProductRepo:         productRepo,
		ProductCategoryRepo: categoryRepo,
	}
}

func (s *ProductService) CreateProduct(product *model.Product) error {
	return s.ProductRepo.Create(product)
}

func (s *ProductService) GetProductByID(id uint) (*model.Product, error) {
	return s.ProductRepo.GetByID(id)
}

func (s *ProductService) GetAllProducts(params map[string]interface{}) ([]model.Product, error) {
	return s.ProductRepo.GetAll(params)
}

func (s *ProductService) UpdateProduct(product *model.Product) error {
	return s.ProductRepo.Update(product)
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.ProductRepo.Delete(id)
}

func (s *ProductService) GetProductsByCreatorID(creatorID uint) ([]model.Product, error) {
	return s.ProductRepo.GetProductsByCreatorID(creatorID)
}

// ProductCategory Service
func (s *ProductService) CreateCategory(category *model.ProductCategory) error {
	return s.ProductCategoryRepo.Create(category)
}

func (s *ProductService) GetCategoryByID(id uint) (*model.ProductCategory, error) {
	return s.ProductCategoryRepo.GetByID(id)
}

func (s *ProductService) GetAllCategories() ([]model.ProductCategory, error) {
	return s.ProductCategoryRepo.GetAll()
}

func (s *ProductService) UpdateCategory(category *model.ProductCategory) error {
	return s.ProductCategoryRepo.Update(category)
}

func (s *ProductService) DeleteCategory(id uint) error {
	return s.ProductCategoryRepo.Delete(id)
}

func (s *ProductService) GetChildCategories(parentID uint) ([]model.ProductCategory, error) {
	return s.ProductCategoryRepo.GetChildCategories(parentID)
}
