package service

import (
	"context"
	"testing"

	"2501YTC/app/cart/biz/dal"
	"2501YTC/app/cart/infra/rpc"
	cart "2501YTC/rpc_gen/kitex_gen/cart"
)

func TestEmptyCart_Run(t *testing.T) {
	dal.Init()
	rpc.InitClient()
	ctx := context.Background()
	s := NewEmptyCartService(ctx)
	// init req and assert value
	req := &cart.EmptyCartReq{
		UserId: 1,
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
