package service

import (
	"context"
	"testing"

	"2501YTC/app/user/biz/dal/mysql"

	"github.com/joho/godotenv"

	user "2501YTC/rpc_gen/kitex_gen/user"
)

func TestLogin_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	mysql.Init()
	ctx := context.Background()
	s := NewLoginService(ctx)
	// init req and assert value

	req := &user.LoginReq{
		Email:    "test@user.com",
		Password: "123456",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
}
