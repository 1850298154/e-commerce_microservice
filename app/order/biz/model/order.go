package model

import (
	"context"
	"database/sql"
	"time"

	orderClient "2501YTC/rpc_gen/kitex_gen/order"

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

// CancelType 取消订单类型
type CancelType string

const (
	CancelTypeUser    CancelType = "user"
	CancelTypeTimeout CancelType = "timeout"
)

// Order 订单信息
type Order struct {
	Base
	OrderId string `gorm:"uniqueIndex;size:256"`

	UserId       uint32
	UserCurrency string
	Consignee    Consignee `gorm:"embedded"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderIdRefer;references:OrderId"`
	OrderState OrderState

	CancelTime sql.NullTime `gorm:"type:datetime"`
	CancelType CancelType
	Version    uint32 `gorm:"version;default:1"` // 乐观锁版本号
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

// GetOrder 获取指定订单详情
func GetOrder(ctx context.Context, db *gorm.DB, orderId string) (order Order, err error) {
	err = db.Where(&Order{OrderId: orderId}).Preload("OrderItems").First(&order).Error
	return
}

// UpdateOrderState 更新订单状态
func UpdateOrderState(ctx context.Context, db *gorm.DB, orderId string, state OrderState) error {
	return db.Model(&Order{}).Where(&Order{OrderId: orderId}).Update("order_state", state).Error
}

// UpdateOrder 更新订单信息
func UpdateOrder(ctx context.Context, db *gorm.DB, orderId string, updates map[string]any) error {
	return db.Model(&Order{}).Where(&Order{OrderId: orderId}).Updates(updates).Error
}

// UpdateOrderItems 更新订单项
func UpdateOrderItems(ctx context.Context, db *gorm.DB, orderId string, items []*orderClient.OrderItem) error {
	// 删除原有订单项
	if err := db.Where("order_id_refer = ?", orderId).Delete(&OrderItem{}).Error; err != nil {
		return err
	}

	// 插入新订单项
	orderItems := make([]OrderItem, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, OrderItem{
			ProductId:    item.Item.ProductId,
			Quantity:     item.Item.Quantity,
			Cost:         item.Cost,
			OrderIdRefer: orderId,
		})
	}
	return db.Create(&orderItems).Error
}

// CancelOrder 取消订单
func CancelOrder(ctx context.Context, db *gorm.DB, orderId string, cancelType CancelType, cancelTime int32) error {
	updates := map[string]any{
		"order_state": OrderStateCanceled,
		"cancel_type": cancelType,
		"cancel_time": time.Unix(int64(cancelTime), 0),
	}
	return db.Model(&Order{}).Where(&Order{OrderId: orderId}).Updates(updates).Error
}

// CancelOrderWithVersion 使用乐观锁取消订单
func CancelOrderWithVersion(ctx context.Context, db *gorm.DB, orderId string, cancelType CancelType, cancelTime int32, version uint32) error {
	updates := map[string]any{
		"order_state": OrderStateCanceled,
		"cancel_type": cancelType,
		"cancel_time": time.Unix(int64(cancelTime), 0),
		"version":     version + 1,
	}
	return db.Model(&Order{}).Where(&Order{OrderId: orderId, Version: version}).Updates(updates).Error
}

// DeleteOrder 删除订单
func DeleteOrder(ctx context.Context, db *gorm.DB, orderId string) error {
	return db.Where(&Order{OrderId: orderId}).Delete(&Order{}).Error
}

// DeleteOrderItemByOrderId 删除订单项
func DeleteOrderItemByOrderId(ctx context.Context, db *gorm.DB, orderId string) error {
	return db.Where("order_id_refer = ?", orderId).Delete(&OrderItem{}).Error
}

// GetOrderVersionAndState 获取订单版本号和状态
func GetOrderVersionAndState(ctx context.Context, db *gorm.DB, orderId string) (uint32, OrderState, error) {
	var order Order
	err := db.Select("version,order_state").Where(&Order{OrderId: orderId}).First(&order).Error
	if err != nil {
		return 0, "", err
	}
	return order.Version, order.OrderState, nil
}
