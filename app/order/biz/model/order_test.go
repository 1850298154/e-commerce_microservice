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
)

func setupTestDB() {
	mysql.Init()
}

func cleanTestDB() {
	db := mysql.DB
	db.Exec("delete from `order_item` where order_id_refer = 'test_order_1'")
	db.Exec("delete from `order` where order_id = 'test_order_1'")
}

func TestCRUDOrder(t *testing.T) {
	setupTestDB()
	defer cleanTestDB()
	orderQuery := model.NewOrderQuery(context.Background(), mysql.DB)

	// Test create order
	order := &model.Order{
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
	err := orderQuery.CreateOrder(order)
	assert.NoError(t, err)

	// Test get order
	getOrder, err := orderQuery.GetOrder("test_order_1")
	assert.NoError(t, err)
	assert.Equal(t, order.OrderId, getOrder.OrderId)

	// Test update order state
	err = orderQuery.UpdateOrderState("test_order_1", model.OrderStatePaid)
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
	err = orderQuery.UpdateOrderItems("test_order_1", items)
	assert.NoError(t, err)

	// Test list orders
	orders, err := orderQuery.ListOrder(1)
	assert.NoError(t, err)
	assert.Len(t, orders, 1)

	// Test cancel order
	err = orderQuery.CancelOrder("test_order_1", model.CancelTypeUser, int32(time.Now().Unix()))
	assert.NoError(t, err)

	// Test get version and state
	version, state, err := orderQuery.GetOrderVersionAndState("test_order_1")
	assert.NoError(t, err)
	assert.Equal(t, model.OrderStateCanceled, state)
	assert.Equal(t, uint32(1), version)

	// Test delete order
	err = orderQuery.DeleteOrderItemByOrderId("test_order_1")
	assert.NoError(t, err)
	err = orderQuery.DeleteOrder("test_order_1")
	assert.NoError(t, err)
}
