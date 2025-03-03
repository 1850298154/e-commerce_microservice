package service

import (
	"context"
	"testing"

	"2501YTC/app/auth/biz/dal/mysql"
	"2501YTC/app/auth/biz/dal/redis"

	"github.com/joho/godotenv"

	auth "2501YTC/rpc_gen/kitex_gen/auth"
)

func TestRenewTokenByRPC_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	mysql.Init()
	redis.Init()
	tests := []struct {
		name    string
		req     *auth.RenewTokenReq
		wantErr bool
	}{
		{
			name: "刷新token正常",
			req: &auth.RenewTokenReq{
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMsIlJvbGUiOjIsImV4cCI6MTczOTgxMDEwMCwianRpIjoiY2M3ZmFlMTktMzRlMi00ZjNhLWJmNjUtMmU0ZmE5YjI0MzgwIiwiaWF0IjoxNzM5MjA1MzAwLCJpc3MiOiJnb21hbGwifQ.92AYrc8XUeTt1ERFFH0kPvM8dZ2hstjcIK9UXQRkA-0",
			},
			wantErr: false,
		},
		{
			name: "refreshtoken未空",
			req: &auth.RenewTokenReq{
				RefreshToken: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewRenewTokenByRPCService(ctx)

			resp, err := s.Run(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("TestRenewTokenByRPC_Run 期望错误但是没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("TestRenewTokenByRPC_Run 错误 = %v", err)
				return
			}

			if resp == nil {
				t.Error("TestRenewTokenByRPC_Run 响应不应该为空")
				return
			}
		})
	}
}
