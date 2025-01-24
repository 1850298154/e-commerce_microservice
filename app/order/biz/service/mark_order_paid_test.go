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

func TestMarkOrderPaid_Run(t *testing.T) {
	dal.Init()
	ctx := context.Background()
	s := NewMarkOrderPaidService(ctx)

	// Test invalid user ID
	invalidUserReq := &order.MarkOrderPaidReq{
		OrderId: "test123",
		UserId:  0,
	}
	resp, err := s.Run(invalidUserReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Test invalid order ID
	invalidOrderReq := &order.MarkOrderPaidReq{
		OrderId: "",
		UserId:  1,
	}
	resp, err = s.Run(invalidOrderReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Test non-existent order
	nonExistentReq := &order.MarkOrderPaidReq{
		OrderId: "non-existent",
		UserId:  999,
	}
	resp, err = s.Run(nonExistentReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Create test order
	testOrder := &model.Order{
		OrderId:      "test123",
		UserId:       1,
		UserCurrency: "USD",
		OrderState:   model.OrderStatePlaced,
		Consignee: model.Consignee{
			Email:   "test@test.com",
			Country: "US",
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
	mysql.DB.Create(testOrder)

	// Test successful mark as paid
	req := &order.MarkOrderPaidReq{
		OrderId: "test123",
		UserId:  1,
	}
	resp, err = s.Run(req)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	// Verify order state updated
	updatedOrder, err := model.GetOrder(ctx, mysql.DB, req.OrderId)
	assert.Nil(t, err)
	assert.Equal(t, model.OrderStatePaid, updatedOrder.OrderState)

	// Clean up
	_ = mysql.DB.Transaction(func(tx *gorm.DB) error {
		assert.NoError(t, model.DeleteOrderItemByOrderId(ctx, tx, "test123"))
		assert.NoError(t, model.DeleteOrder(ctx, tx, "test123"))
		return nil
	})
}
