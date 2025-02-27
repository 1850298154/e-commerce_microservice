package service

import (
	"context"
	"testing"

	"2501YTC/app/ai/infra/rpc"

	"github.com/joho/godotenv"

	ai "2501YTC/rpc_gen/kitex_gen/ai"
)

func TestAutoOrder_Run(t *testing.T) {
	rpc.InitClient()
	_ = godotenv.Load("../../.env")
	ctx := context.Background()
	s := NewAutoOrderService(ctx)
	// init req and assert value
	req := &ai.AutoOrderReq{
		UserId:  1,
		Content: "购买2件衬衫和2个小米手环9。",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
