package service

import (
	"context"
	"testing"

	auth "2501YTC/rpc_gen/kitex_gen/auth"
)

func TestDeleteTokenByRPC_Run(t *testing.T) {
	ctx := context.Background()
	s := NewDeleteTokenByRPCService(ctx)
	// init req and assert value
	req := &auth.DeleteTokenReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
	// todo: edit your unit test
}
