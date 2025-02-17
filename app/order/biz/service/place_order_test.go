package service

import (
	"context"
	"testing"
	"time"

	"2501YTC/app/order/biz/dal"
	"2501YTC/app/order/biz/dal/mq"
	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	"2501YTC/rpc_gen/kitex_gen/cart"
	"2501YTC/rpc_gen/kitex_gen/order"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPlaceOrder(t *testing.T) {
	dal.Init()
	mq.ProducerInstance, _ = mq.NewProducer(3) // 3s timeout
	consumer, _ := mq.NewConsumer(mysql.DB)

	ready := make(chan struct{})
	go func() {
		close(ready)
		_ = consumer.Consume()
	}()
	<-ready
	ctx := context.Background()
	s := NewPlaceOrderService(ctx)

	tests := []struct {
		name    string
		req     *order.PlaceOrderReq
		wantErr bool
	}{
		{
			name: "empty order items",
			req: &order.PlaceOrderReq{
				UserId:       1,
				UserCurrency: "USD",
				Email:        "test@test.com",
				OrderItems:   []*order.OrderItem{},
			},
			wantErr: true,
		},
		{
			name: "empty user id",
			req: &order.PlaceOrderReq{
				UserId: 0,
				OrderItems: []*order.OrderItem{
					{
						Cost: 100,
						Item: &cart.CartItem{
							ProductId: 1,
							Quantity:  1,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "multiple items",
			req: &order.PlaceOrderReq{
				UserId:       1,
				UserCurrency: "CNY",
				Email:        "test@test.com",
				Address: &order.Address{
					Country:       "China",
					State:         "Beijing",
					City:          "Beijing",
					StreetAddress: "123 Test St",
				},
				OrderItems: []*order.OrderItem{
					{
						Cost: 100,
						Item: &cart.CartItem{
							ProductId: 1,
							Quantity:  2,
						},
					},
					{
						Cost: 200,
						Item: &cart.CartItem{
							ProductId: 2,
							Quantity:  1,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	db := mysql.DB
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.Run(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				// clean up
				_ = mysql.DB.Transaction(func(tx *gorm.DB) error {
					cleanQuery := model.NewOrderQuery(ctx, tx)
					assert.NoError(t, cleanQuery.DeleteOrderItemByOrderId(resp.Order.OrderId))
					assert.NoError(t, cleanQuery.DeleteOrder(resp.Order.OrderId))
					return nil
				})
				db.Exec("delete from `order_item` where order_id_refer = ?", resp.Order.OrderId)
				db.Exec("delete from `order` where order_id = ?", resp.Order.OrderId)
			}
		})
	}
	mq.ProducerInstance.Stop()
	consumer.Stop()
}

func TestCancelOrderWithTimeout(t *testing.T) {
	dal.Init()
	mq.ProducerInstance, _ = mq.NewProducer(3) // 3s timeout
	consumer, _ := mq.NewConsumer(mysql.DB)

	ready := make(chan struct{})
	go func() {
		close(ready)
		_ = consumer.Consume()
	}()
	<-ready

	ctx := context.Background()
	s := NewPlaceOrderService(ctx)

	// create order
	req := &order.PlaceOrderReq{
		UserId:       1,
		UserCurrency: "CNY",
		Email:        "test@test.com",
		Address: &order.Address{
			Country:       "China",
			State:         "Beijing",
			City:          "Beijing",
			StreetAddress: "123 Test St",
		},
		OrderItems: []*order.OrderItem{
			{
				Cost: 100,
				Item: &cart.CartItem{
					ProductId: 1,
					Quantity:  2,
				},
			},
			{
				Cost: 200,
				Item: &cart.CartItem{
					ProductId: 2,
					Quantity:  1,
				},
			},
		},
	}
	resp, err := s.Run(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// cancel order when timeout
	time.Sleep(10 * time.Second)
	// check order state
	orderQuery := model.NewOrderQuery(ctx, mysql.DB)
	canceledOrder, err := orderQuery.GetOrder(resp.Order.OrderId)
	assert.Equal(t, model.OrderStateCanceled, canceledOrder.OrderState)
	assert.Equal(t, model.CancelTypeTimeout, canceledOrder.CancelType)

	// clean up
	_ = mysql.DB.Transaction(func(tx *gorm.DB) error {
		cleanQuery := model.NewOrderQuery(ctx, tx)
		assert.NoError(t, cleanQuery.DeleteOrderItemByOrderId(resp.Order.OrderId))
		assert.NoError(t, cleanQuery.DeleteOrder(resp.Order.OrderId))
		return nil
	})
	time.Sleep(3 * time.Second)
	db := mysql.DB
	db.Exec("delete from `order_item` where order_id_refer = ?", resp.Order.OrderId)
	db.Exec("delete from `order` where order_id = ?", resp.Order.OrderId)
	mq.ProducerInstance.Stop()
	consumer.Stop()
}
