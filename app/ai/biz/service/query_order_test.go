package service

import (
	"context"
	"testing"

	"2501YTC/app/ai/infra/rpc"

	"github.com/joho/godotenv"

	ai "2501YTC/rpc_gen/kitex_gen/ai"
)

func TestQueryOrder_Run(t *testing.T) {
	rpc.InitClient()
	_ = godotenv.Load("../../.env")
	ctx := context.Background()
	s := NewQueryOrderService(ctx)
	// init req and assert value
	req := &ai.OrderQueryReq{
		UserId:  1,
		Content: "查找六天前的订单",
	}

	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
