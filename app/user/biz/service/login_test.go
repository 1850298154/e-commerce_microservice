package service

import (
	"context"
	"testing"

	user "2501YTC/rpc_gen/kitex_gen/user"
)

func TestLogin_Run(t *testing.T) {
	ctx := context.Background()
	s := NewLoginService(ctx)
	// init req and assert value

	req := &user.LoginReq{}
	resp, err := s.Run(req)
	t.Logf("apiErr: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
