package service

import (
	"context"
	"testing"

	cart "2501YTC/rpc_gen/kitex_gen/cart"
)

func TestGetCart_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetCartService(ctx)
	// init req and assert value

	req := &cart.GetCartReq{}
	resp, err := s.Run(req)
	t.Logf("apiErr: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
