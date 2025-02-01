package service

import (
	auth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"
	"testing"
)

func TestRenewTokenByRPC_Run(t *testing.T) {
	ctx := context.Background()
	s := NewRenewTokenByRPCService(ctx)
	// init req and assert value

	req := &auth.RenewTokenReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
	// todo: edit your unit test
}
