package model_test

import (
	"context"
	"testing"
	"time"

	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	"2501YTC/rpc_gen/kitex_gen/cart"
	orderClient "2501YTC/rpc_gen/kitex_gen/order"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	mysql.Init()
	return mysql.DB
}

func TestCRUDOrder(t *testing.T) {
	db := setupTestDB()
	ctx := context.Background()

	// Test create order
	order := model.Order{
		OrderId:      "test_order_1",
		UserId:       1,
		UserCurrency: "USD",
		Consignee: model.Consignee{
			Email:         "test@test.com",
			RecipientName: "Test User",
			PhoneNumber:   "1234567890",
		},
		OrderState: model.OrderStatePlaced,
	}
	err := db.Create(&order).Error
	assert.NoError(t, err)

	// Test get order
	getOrder, err := model.GetOrder(ctx, db, "test_order_1")
	assert.NoError(t, err)
	assert.Equal(t, order.OrderId, getOrder.OrderId)

	// Test update order state
	err = model.UpdateOrderState(ctx, db, "test_order_1", model.OrderStatePaid)
	assert.NoError(t, err)

	// Test update order items
	items := []*orderClient.OrderItem{
		{
			Item: &cart.CartItem{
				ProductId: 1,
				Quantity:  2,
			},
			Cost: 20.0,
		},
	}
	err = model.UpdateOrderItems(ctx, db, "test_order_1", items)
	assert.NoError(t, err)

	// Test list orders
	orders, err := model.ListOrder(ctx, db, 1)
	assert.NoError(t, err)
	assert.Len(t, orders, 1)

	// Test cancel order
	err = model.CancelOrder(ctx, db, "test_order_1", model.CancelTypeUser, int32(time.Now().Unix()))
	assert.NoError(t, err)

	// Test get version and state
	version, state, err := model.GetOrderVersionAndState(ctx, db, "test_order_1")
	assert.NoError(t, err)
	assert.Equal(t, model.OrderStateCanceled, state)
	assert.Equal(t, uint32(1), version)

	// Test delete order
	err = model.DeleteOrderItemByOrderId(ctx, db, "test_order_1")
	assert.NoError(t, err)
	err = model.DeleteOrder(ctx, db, "test_order_1")
	assert.NoError(t, err)
}
