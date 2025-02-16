package service

import (
	"context"
	"testing"
	ai "2501YTC/rpc_gen/kitex_gen/ai"
)

func TestAutoOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewAutoOrderService(ctx)
	// init req and assert value

	req := &ai.AutoOrderReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
