package service

import (
	"context"
	"testing"
	"time"

	"2501YTC/app/order/biz/dal"
	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	"2501YTC/rpc_gen/kitex_gen/order"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCancelOrder_Run(t *testing.T) {
	dal.Init()
	ctx := context.Background()
	s := NewCancelOrderService(ctx)

	// Test invalid user ID
	invalidUserReq := &order.CancelOrderReq{
		OrderId: "test123",
		UserId:  0,
	}
	resp, err := s.Run(invalidUserReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Test invalid order ID
	invalidOrderReq := &order.CancelOrderReq{
		OrderId: "",
		UserId:  1,
	}
	resp, err = s.Run(invalidOrderReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Test non-existent order
	nonExistentReq := &order.CancelOrderReq{
		OrderId: "non-existent",
		UserId:  999,
	}
	resp, err = s.Run(nonExistentReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Create test order
	testOrder1 := &model.Order{
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
	orderQuery := model.NewOrderQuery(ctx, mysql.DB)

	assert.NoError(t, orderQuery.CreateOrder(testOrder1))

	// Test user cancel
	userCancelReq := &order.CancelOrderReq{
		OrderId:     "test123",
		UserId:      1,
		TimedCancel: false,
		CancelTime:  int32(time.Now().Unix()),
	}
	resp, err = s.Run(userCancelReq)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// Verify order state updated
	updatedOrder, err := orderQuery.GetOrder(userCancelReq.OrderId)
	assert.Nil(t, err)
	assert.Equal(t, model.OrderStateCanceled, updatedOrder.OrderState)

	_ = mysql.DB.Transaction(func(tx *gorm.DB) error {
		cleanQuery := model.NewOrderQuery(ctx, tx)
		assert.NoError(t, cleanQuery.DeleteOrderItemByOrderId("test123"))
		assert.NoError(t, cleanQuery.DeleteOrder("test123"))
		return nil
	})

	db := mysql.DB
	db.Exec("delete from `order_item` where order_id_refer = 'test123'")
	db.Exec("delete from `order` where order_id = 'test123'")
}
