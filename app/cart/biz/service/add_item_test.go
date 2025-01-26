package service

import (
	"context"
	"testing"

	cart "2501YTC/rpc_gen/kitex_gen/cart"
)

func TestAddItem_Run(t *testing.T) {
	ctx := context.Background()
	s := NewAddItemService(ctx)
	// init req and assert value
	req := &cart.AddItemReq{
		UserId: 1,
		Item: &cart.CartItem{
			ProductId: 1,
			Quantity:  1,
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
