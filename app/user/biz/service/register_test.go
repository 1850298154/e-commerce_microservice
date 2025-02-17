package service

import (
	"context"
	"testing"

	"2501YTC/app/user/biz/dal/mysql"

	"github.com/joho/godotenv"

	user "2501YTC/rpc_gen/kitex_gen/user"
)

func TestRegister_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	mysql.Init()
	ctx := context.Background()
	s := NewRegisterService(ctx)
	// init req and assert value

	req := &user.RegisterReq{
		Email:           "test@admin.com",
		Password:        "123456",
		ConfirmPassword: "123456",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
}
