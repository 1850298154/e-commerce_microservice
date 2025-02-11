package service

import (
	"context"
	"testing"

	"2501YTC/app/user/biz/dal/mysql"

	"github.com/joho/godotenv"

	user "2501YTC/rpc_gen/kitex_gen/user"
)

func TestUpdateUserRole_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	mysql.Init()
	ctx := context.Background()
	s := NewUpdateUserRoleService(ctx)
	// init req and assert value

	req := &user.UpdateUserRoleReq{
		UserId: 2,
		Role:   0,
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
}
