package service

import (
	"context"
	"testing"

	"2501YTC/app/auth/biz/dal/mysql"
	"2501YTC/app/auth/biz/dal/redis"

	"github.com/joho/godotenv"

	auth "2501YTC/rpc_gen/kitex_gen/auth"
)

func TestDeleteTokenByRPC_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	mysql.Init()
	redis.Init()
	tests := []struct {
		name    string
		req     *auth.DeleteTokenReq
		wantErr bool
	}{
		{
			name: "删除token正常",
			req: &auth.DeleteTokenReq{
				Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMsIlJvbGUiOjIsImV4cCI6MTczOTIwODkwMCwianRpIjoiY2M3ZmFlMTktMzRlMi00ZjNhLWJmNjUtMmU0ZmE5YjI0MzgwIiwiaWF0IjoxNzM5MjA1MzAwLCJpc3MiOiJnb21hbGwifQ.Sxgk_n6WaL0I7BXIUQHwAW_DEiDrPB9mEAYFSIiG_dU",
			},
			wantErr: false,
		},
		{
			name: "token未空",
			req: &auth.DeleteTokenReq{
				Token: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewDeleteTokenByRPCService(ctx)

			resp, err := s.Run(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("TestDeleteTokenByRPC_Run 期望错误但是没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("TestDeleteTokenByRPC_Run 错误 = %v", err)
				return
			}

			if resp == nil {
				t.Error("TestDeleteTokenByRPC_Run 响应不应该为空")
				return
			}

		})
	}
}
