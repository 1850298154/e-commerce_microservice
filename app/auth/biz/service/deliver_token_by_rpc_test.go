package service

import (
	"context"
	"testing"

	"2501YTC/app/auth/biz/dal/mysql"
	"2501YTC/app/auth/biz/dal/redis"

	"github.com/joho/godotenv"

	auth "2501YTC/rpc_gen/kitex_gen/auth"
)

func TestDeliverTokenByRPC_Run(t *testing.T) {

	// token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMsIlJvbGUiOjIsImV4cCI6MTczOTIwODkwMCwianRpIjoiY2M3ZmFlMTktMzRlMi00ZjNhLWJmNjUtMmU0ZmE5YjI0MzgwIiwiaWF0IjoxNzM5MjA1MzAwLCJpc3MiOiJnb21hbGwifQ.Sxgk_n6WaL0I7BXIUQHwAW_DEiDrPB9mEAYFSIiG_dU"
	// refresh_token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMsIlJvbGUiOjIsImV4cCI6MTczOTgxMDEwMCwianRpIjoiY2M3ZmFlMTktMzRlMi00ZjNhLWJmNjUtMmU0ZmE5YjI0MzgwIiwiaWF0IjoxNzM5MjA1MzAwLCJpc3MiOiJnb21hbGwifQ.92AYrc8XUeTt1ERFFH0kPvM8dZ2hstjcIK9UXQRkA-0"
	// token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMsIlJvbGUiOjIsImV4cCI6MTczOTIwODkwMCwianRpIjoiY2M3ZmFlMTktMzRlMi00ZjNhLWJmNjUtMmU0ZmE5YjI0MzgwIiwiaWF0IjoxNzM5MjA1MzAwLCJpc3MiOiJnb21hbGwifQ.Sxgk_n6WaL0I7BXIUQHwAW_DEiDrPB9mEAYFSIiG_dU"
	// refresh_token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMsIlJvbGUiOjIsImV4cCI6MTczOTgxMDEwMCwianRpIjoiY2M3ZmFlMTktMzRlMi00ZjNhLWJmNjUtMmU0ZmE5YjI0MzgwIiwiaWF0IjoxNzM5MjA1MzAwLCJpc3MiOiJnb21hbGwifQ.92AYrc8XUeTt1ERFFH0kPvM8dZ2hstjcIK9UXQRkA-0"
	_ = godotenv.Load("../../.env")
	mysql.Init()
	redis.Init()
	tests := []struct {
		name    string
		req     *auth.DeliverTokenReq
		wantErr bool
	}{
		{
			name: "分发token正常",
			req: &auth.DeliverTokenReq{
				UserId: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewDeliverTokenByRPCService(ctx)

			resp, err := s.Run(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("TestDeliverTokenByRPC_Run 期望错误但是没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("TestDeliverTokenByRPC_Run 错误 = %v", err)
				return
			}

			if resp == nil {
				t.Error("TestDeliverTokenByRPC_Run 响应不应该为空")
				return
			}

		})
	}
}
