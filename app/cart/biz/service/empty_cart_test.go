package service

import (
	"context"
	"testing"

	cart "2501YTC/rpc_gen/kitex_gen/cart"
)

func TestEmptyCart_Run(t *testing.T) {
	ctx := context.Background()
	s := NewEmptyCartService(ctx)
	// init req and assert value

	req := &cart.EmptyCartReq{}
	resp, err := s.Run(req)
	t.Logf("apiErr: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
