package service

import (
	"2501YTC/app/auth/biz/dal/mysql"
	"2501YTC/app/auth/biz/dal/redis"
	"context"
	"github.com/joho/godotenv"
	"testing"

	auth "2501YTC/rpc_gen/kitex_gen/auth"
)

func TestRenewTokenByRPC_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	mysql.Init()
	redis.Init()
	ctx := context.Background()
	s := NewRenewTokenByRPCService(ctx)
	// init req and assert value

	req := &auth.RenewTokenReq{
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjQsIlJvbGUiOjQsImV4cCI6MTczOTgxMDU3OCwianRpIjoiZjkxYWE4MDctZTc1ZS00Y2Y0LThhOTktNDUwMmU4NTU1YzM4IiwiaWF0IjoxNzM5MjA1Nzc4LCJpc3MiOiJnb21hbGwifQ.lNTLA8pvDXKKEYzeBTiJZdeKD6GEb9VtgDWKqMbswnU",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
	// todo: edit your unit test
}
