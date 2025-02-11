package service

import (
	"2501YTC/app/auth/biz/dal/mysql"
	"2501YTC/app/auth/biz/dal/redis"
	auth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"
	"github.com/joho/godotenv"
	"testing"
)

func TestDeleteTokenByRPC_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	mysql.Init()
	redis.Init()
	ctx := context.Background()
	s := NewDeleteTokenByRPCService(ctx)
	// init req and assert value
	req := &auth.DeleteTokenReq{
		Token: "revoked",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
	// todo: edit your unit test
}
