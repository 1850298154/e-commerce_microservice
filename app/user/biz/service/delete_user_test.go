package service

import (
	"context"
	"testing"

	"2501YTC/app/user/biz/dal/mysql"

	user "2501YTC/rpc_gen/kitex_gen/user"
)

func TestDeleteUser_Run(t *testing.T) {
	//_ = godotenv.Load("../../.env")
	mysql.Init()
	ctx := context.Background()
	s := NewDeleteUserService(ctx)
	// init req and assert value

	req := &user.DeleteUserReq{
		UserId: 2,
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
}
