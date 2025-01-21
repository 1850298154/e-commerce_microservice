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

func TestUpdateOrderScenarios(t *testing.T) {
	dal.Init()
	ctx := context.Background()
	s := NewUpdateOrderService(ctx)

	// Test invalid user ID
	invalidUserReq := &order.UpdateOrderReq{
		OrderId: "test123",
		UserId:  0,
	}
	resp, err := s.Run(invalidUserReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Test invalid order ID
	invalidOrderReq := &order.UpdateOrderReq{
		OrderId: "",
		UserId:  1,
	}
	resp, err = s.Run(invalidOrderReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Test updating non-existent order
	nonExistentReq := &order.UpdateOrderReq{
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
		Consignee: model.Consignee{
			Email:         "old@test.com",
			Country:       "US",
			StreetAddress: "Old Street",
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

	// Test updating only email
	emailUpdateReq := &order.UpdateOrderReq{
		OrderId:  "test123",
		UserId:   1,
		NewEmail: "new@test.com",
	}
	resp, err = s.Run(emailUpdateReq)
	assert.Nil(t, err)
	assert.True(t, resp.Success)

	updatedOrder, err := model.GetOrder(ctx, mysql.DB, emailUpdateReq.UserId, emailUpdateReq.OrderId)
	assert.Nil(t, err)
	assert.Equal(t, emailUpdateReq.NewEmail, updatedOrder.Consignee.Email)

	// Test updating full address
	addressUpdateReq := &order.UpdateOrderReq{
		OrderId: "test123",
		UserId:  1,
		NewAddress: &order.Address{
			City:          "New City",
			State:         "New State",
			Country:       "New Country",
			ZipCode:       0,
			StreetAddress: "New Street",
		},
	}
	resp, err = s.Run(addressUpdateReq)
	assert.Nil(t, err)
	assert.True(t, resp.Success)

	updatedOrder, err = model.GetOrder(ctx, mysql.DB, addressUpdateReq.UserId, addressUpdateReq.OrderId)
	assert.Nil(t, err)
	assert.Equal(t, addressUpdateReq.NewAddress.City, updatedOrder.Consignee.City)
	assert.Equal(t, addressUpdateReq.NewAddress.State, updatedOrder.Consignee.State)
	assert.Equal(t, addressUpdateReq.NewAddress.Country, updatedOrder.Consignee.Country)
	assert.Equal(t, addressUpdateReq.NewAddress.ZipCode, updatedOrder.Consignee.ZipCode)
	assert.Equal(t, addressUpdateReq.NewAddress.StreetAddress, updatedOrder.Consignee.StreetAddress)

	// Clean up
	_ = mysql.DB.Transaction(func(tx *gorm.DB) error {
		assert.NoError(t, model.DeleteOrderItemByOrderId(ctx, tx, "test123"))
		assert.NoError(t, model.DeleteOrder(ctx, tx, 1, "test123"))
		return nil
	})
}
