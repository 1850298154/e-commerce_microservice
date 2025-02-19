package service

import (
	"context"
	"testing"

	"2501YTC/app/order/biz/dal"
	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	order "2501YTC/rpc_gen/kitex_gen/order"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetOrder_Run(t *testing.T) {
	dal.Init()
	ctx := context.Background()
	s := NewGetOrderService(ctx)

	defer func() {
		// Clean up
		_ = mysql.DB.Transaction(func(tx *gorm.DB) error {
			cleanQuery := model.NewOrderQuery(ctx, tx)
			assert.NoError(t, cleanQuery.DeleteOrderItemByOrderId("test789"))
			assert.NoError(t, cleanQuery.DeleteOrder("test789"))
			return nil
		})
		db := mysql.DB
		db.Exec("delete from `order_item` where order_id_refer = 'test789'")
		db.Exec("delete from `order` where order_id = 'test789'")
	}()

	// Test invalid order ID
	invalidReq := &order.GetOrderReq{
		OrderId: "",
	}
	resp, err := s.Run(invalidReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Create test order
	testOrder := &model.Order{
		OrderId:      "test789",
		UserId:       1,
		UserCurrency: "USD",
		Consignee: model.Consignee{
			Email:         "test@example.com",
			Country:       "US",
			StreetAddress: "123 Test Street",
		},
		OrderItems: []model.OrderItem{
			{
				ProductId:    1,
				OrderIdRefer: "test789",
				Quantity:     2,
				Cost:         15.99,
			},
		},
	}

	// Insert test data
	orderQuery := model.NewOrderQuery(ctx, mysql.DB)
	assert.NoError(t, orderQuery.CreateOrder(testOrder))

	// Test getting existing order
	getReq := &order.GetOrderReq{
		UserId:  1,
		OrderId: "test789",
	}
	resp, err = s.Run(getReq)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Order)

	// Verify order details
	assert.Equal(t, "test789", resp.Order.Order.OrderId)
	assert.Equal(t, uint32(1), resp.Order.Order.UserId)
	assert.Equal(t, "USD", resp.Order.Order.UserCurrency)
	assert.Equal(t, "test@example.com", resp.Order.Order.Email)
	assert.Equal(t, "US", resp.Order.Order.Address.Country)
	assert.Equal(t, "123 Test Street", resp.Order.Order.Address.StreetAddress)
	assert.Equal(t, 1, len(resp.Order.Order.OrderItems))
	assert.Equal(t, uint32(1), resp.Order.Order.OrderItems[0].Item.ProductId)
	assert.Equal(t, int32(2), resp.Order.Order.OrderItems[0].Item.Quantity)
	assert.Equal(t, float32(15.99), resp.Order.Order.OrderItems[0].Cost)

	// Test getting non-existent order
	nonExistentReq := &order.GetOrderReq{
		OrderId: "nonexistent",
	}
	resp, err = s.Run(nonExistentReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}
