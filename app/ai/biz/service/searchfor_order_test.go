package service

import (
	ai "2501YTC/rpc_gen/kitex_gen/ai"
	"context"
	"testing"
)

func TestSearchforOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewSearchforOrderService(ctx)
	// init req and assert value

	req := &ai.SearchforOrderReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
	// todo: edit your unit test
}
