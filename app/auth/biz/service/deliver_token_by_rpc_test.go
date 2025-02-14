package service

import (
	"2501YTC/app/auth/biz/dal/mysql"
	"2501YTC/app/auth/biz/dal/redis"
	"2501YTC/app/auth/conf"
	"context"
	"fmt"
	"testing"

	"github.com/joho/godotenv"

	auth "2501YTC/rpc_gen/kitex_gen/auth"
)

func TestDeliverTokenByRPC_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	mysql.Init()
	redis.Init()
	fmt.Println(conf.GetConf().MySQL.DSN)
	ctx := context.Background()
	s := NewDeliverTokenByRPCService(ctx)
	// init req and assert value
	req := &auth.DeliverTokenReq{
		UserId: 123456,
		Role:   1,
	}
	// token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMsIlJvbGUiOjIsImV4cCI6MTczOTIwODkwMCwianRpIjoiY2M3ZmFlMTktMzRlMi00ZjNhLWJmNjUtMmU0ZmE5YjI0MzgwIiwiaWF0IjoxNzM5MjA1MzAwLCJpc3MiOiJnb21hbGwifQ.Sxgk_n6WaL0I7BXIUQHwAW_DEiDrPB9mEAYFSIiG_dU"
	// refresh_token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMsIlJvbGUiOjIsImV4cCI6MTczOTgxMDEwMCwianRpIjoiY2M3ZmFlMTktMzRlMi00ZjNhLWJmNjUtMmU0ZmE5YjI0MzgwIiwiaWF0IjoxNzM5MjA1MzAwLCJpc3MiOiJnb21hbGwifQ.92AYrc8XUeTt1ERFFH0kPvM8dZ2hstjcIK9UXQRkA-0"
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
	// todo: edit your unit test
}
