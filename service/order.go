package service

import (
	"myApp/model"
	"myApp/repository"
)

type OrderService interface {
	CreateOrder(order *model.Order) (*model.Order, error)
	GetOrder(id uint) (*model.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(order *model.Order) (*model.Order, error) {
	return s.repo.CreateOrder(order)
}

func (s *orderService) GetOrder(id uint) (*model.Order, error) {
	return s.repo.GetOrder(id)
}
