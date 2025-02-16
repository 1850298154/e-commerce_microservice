package service

import (
	"context"
	"testing"
	ai "2501YTC/rpc_gen/kitex_gen/ai"
)

func TestQueryOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewQueryOrderService(ctx)
	// init req and assert value

	req := &ai.OrderQueryReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
