package repository

import (
	"myApp/model"
)

type OrderRepository interface {
	CreateOrder(order *model.Order) (*model.Order, error)
	GetOrder(id uint) (*model.Order, error)
}

type orderRepository struct{}

func NewOrderRepository() OrderRepository {
	return &orderRepository{}
}

func (r *orderRepository) CreateOrder(order *model.Order) (*model.Order, error) {
	// 实现数据库操作
	err := model.GetDB().Create(order).Error
	return order, err
}

func (r *orderRepository) GetOrder(id uint) (*model.Order, error) {
	var order model.Order
	err := model.GetDB().First(&order, id).Error
	return &order, err
}
