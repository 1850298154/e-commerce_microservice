package service

import (
	"context"
	"testing"
	order "2501YTC/rpc_gen/kitex_gen/order"
)
// TODO
func TestCancelOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCancelOrderService(ctx)
	// init req and assert value

	req := &order.CancelOrderReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
