package model

type Order struct {
	BaseModel
	UserID     uint   
	ProductID  uint   
	Quantity   int    
	TotalPrice float64
	Status     string 
}

// 订单状态常量
const (
	OrderPending   = "pending"
	OrderPaid      = "paid"
	OrderShipped   = "shipped"
	OrderCompleted = "completed"
	OrderCancelled = "cancelled"
)