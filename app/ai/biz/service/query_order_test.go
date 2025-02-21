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
	// orderDetail := "{\\\"order_items\\\":[{\\\"item\\\":{\\\"product_id\\\":1,\\\"quantity\\\":2},\\\"cost\\\":100},{\\\"item\\\":{\\\"product_id\\\":2,\\\"quantity\\\":1},\\\"cost\\\":200}],\\\"order_id\\\":\\\"5051b674-ec0a-11ef-938c-526b95ad440b\\\",\\\"user_id\\\":1,\\\"user_currency\\\":\\\"CNY\\\",\\\"address\\\":{\\\"street_address\\\":\\\"123 Test St\\\",\\\"city\\\":\\\"Beijing\\\",\\\"country\\\":\\\"China\\\"},\\\"email\\\":\\\"test@test.com\\\",\\\"created_at\\\":1739671451}"
	req := &ai.OrderQueryReq{
		UserId:  1,
		Content: "查找购买过衬衫订单",
	}

	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
