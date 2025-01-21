package service

import (
	"context"
	"testing"

	"2501YTC/app/order/biz/dal"
	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	"2501YTC/rpc_gen/kitex_gen/order"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestListOrder_Run(t *testing.T) {
	dal.Init()
	ctx := context.Background()
	s := NewListOrderService(ctx)

	// Test invalid user ID
	invalidUserReq := &order.ListOrderReq{
		UserId: 0,
	}
	resp, err := s.Run(invalidUserReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Create test orders
	testOrder1 := &model.Order{
		OrderId:      "test123",
		UserId:       1,
		UserCurrency: "USD",
		Consignee: model.Consignee{
			Email:         "test1@test.com",
			Country:       "US",
			StreetAddress: "Street 1",
		},
		OrderItems: []model.OrderItem{
			{
				ProductId:    1,
				OrderIdRefer: "test123",
				Quantity:     1,
				Cost:         10,
			},
		},
	}

	testOrder2 := &model.Order{
		OrderId:      "test456",
		UserId:       1,
		UserCurrency: "EUR",
		Consignee: model.Consignee{
			Email:         "test2@test.com",
			Country:       "FR",
			StreetAddress: "Street 2",
		},
		OrderItems: []model.OrderItem{
			{
				ProductId:    2,
				OrderIdRefer: "test456",
				Quantity:     2,
				Cost:         20,
			},
		},
	}

	// Insert test data
	mysql.DB.Create(testOrder1)
	mysql.DB.Create(testOrder2)

	// Test listing orders for user
	listReq := &order.ListOrderReq{
		UserId: 1,
	}
	resp, err = s.Run(listReq)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Orders))

	// Verify first order details
	assert.Equal(t, "test123", resp.Orders[0].OrderId)
	assert.Equal(t, "USD", resp.Orders[0].UserCurrency)
	assert.Equal(t, "test1@test.com", resp.Orders[0].Email)
	assert.Equal(t, "US", resp.Orders[0].Address.Country)
	assert.Equal(t, 1, len(resp.Orders[0].OrderItems))
	assert.Equal(t, float32(10), resp.Orders[0].OrderItems[0].Cost)

	// Test listing orders for non-existent user
	nonExistentReq := &order.ListOrderReq{
		UserId: 999,
	}
	resp, err = s.Run(nonExistentReq)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 0, len(resp.Orders))

	// Clean up
	_ = mysql.DB.Transaction(func(tx *gorm.DB) error {
		assert.NoError(t, model.DeleteOrderItemByOrderId(ctx, tx, "test123"))
		assert.NoError(t, model.DeleteOrderItemByOrderId(ctx, tx, "test456"))
		assert.NoError(t, model.DeleteOrder(ctx, tx, 1, "test123"))
		assert.NoError(t, model.DeleteOrder(ctx, tx, 1, "test456"))
		return nil
	})
}
