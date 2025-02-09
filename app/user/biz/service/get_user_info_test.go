package service

import (
	"context"
	"testing"

	"2501YTC/app/user/biz/dal/mysql"

	"github.com/joho/godotenv"

	user "2501YTC/rpc_gen/kitex_gen/user"
)

func TestGetUserInfo_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	mysql.Init()
	ctx := context.Background()
	s := NewGetUserInfoService(ctx)
	// init req and assert value

	req := &user.GetUserInfoReq{
		UserId: 2,
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
}
