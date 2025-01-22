package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Base 必须包含的字段
type Base struct {
	ID        int `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OrderItem 订单项
type OrderItem struct {
	Base
	ProductId    uint32
	OrderIdRefer string `gorm:"size:256;index"`
	Quantity     int32
	Cost         float32
}

// Consignee 收货人信息
type Consignee struct {
	Email         string
	RecipientName string
	PhoneNumber   string

	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
}

// OrderState 订单状态
type OrderState string

const (
	OrderStatePlaced   OrderState = "placed"   // 订单已下单
	OrderStatePaid     OrderState = "paid"     // 订单已支付
	OrderStateCanceled OrderState = "canceled" // 订单已取消
)

// Order 订单信息
type Order struct {
	Base
	OrderId string `gorm:"uniqueIndex;size:256"`

	UserId       uint32
	UserCurrency string
	Consignee    Consignee `gorm:"embedded"`

	OrderItems  []OrderItem `gorm:"foreignKey:OrderIdRefer;references:OrderId"`
	OrderState  OrderState
	TotalAmount float32
}

// TableName 返回表名
func (oi OrderItem) TableName() string {
	return "order_item"
}

// TableName 返回表名
func (o Order) TableName() string {
	return "order"
}

// ListOrder 获取用户的订单列表
func ListOrder(ctx context.Context, db *gorm.DB, userId uint32) (orders []Order, err error) {
	err = db.Model(&Order{}).Where(&Order{UserId: userId}).Preload("OrderItems").Find(&orders).Error
	return
}

// GetOrder 获取订单详情
func GetOrder(ctx context.Context, db *gorm.DB, userId uint32, orderId string) (order Order, err error) {
	err = db.Where(&Order{UserId: userId, OrderId: orderId}).First(&order).Error
	return
}

// UpdateOrderState 更新订单状态
func UpdateOrderState(ctx context.Context, db *gorm.DB, userId uint32, orderId string, state OrderState) error {
	return db.Model(&Order{}).Where(&Order{UserId: userId, OrderId: orderId}).Update("order_state", state).Error
}
