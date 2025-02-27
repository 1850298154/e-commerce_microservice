package service

import (
	"context"
	"testing"

	"2501YTC/app/auth/biz/dal/mysql"
	"2501YTC/app/auth/biz/dal/redis"

	"github.com/joho/godotenv"

	auth "2501YTC/rpc_gen/kitex_gen/auth"
)

func TestVerifyTokenByRPC_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	mysql.Init()
	redis.Init()
	tests := []struct {
		name    string
		req     *auth.VerifyTokenReq
		wantErr bool
	}{
		{
			name: "验证token正常",
			req: &auth.VerifyTokenReq{
				Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEyMzQ1LCJSb2xlIjoxLCJpc3MiOiJnb21hbGwiLCJleHAiOjE3NDAxMzczMjMsImlhdCI6MTc0MDEzMzcyMywianRpIjoiNDMyOGJmYTgtODFkYS00ODM1LWFiNjQtZmIxZjJkMGQ3NDQwIn0.i5VBBTnRdLcWkTEwOO4uZwrXebyZuakXV3N9CchYmhU",
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEyMzQ1LCJSb2xlIjoxLCJpc3MiOiJnb21hbGwiLCJleHAiOjE3NDA3Mzg1MjMsImlhdCI6MTc0MDEzMzcyMywianRpIjoiNDMyOGJmYTgtODFkYS00ODM1LWFiNjQtZmIxZjJkMGQ3NDQwIn0.o6Gd87mxPStgnIhseNFrHWIM4hAmiF2lZAXlQ1QiXXo",
			},
			wantErr: false,
		},
		{
			name: "refreshtoken或token未空",
			req: &auth.VerifyTokenReq{
				Token:        "",
				RefreshToken: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewVerifyTokenByRPCService(ctx)

			resp, err := s.Run(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("TestVerifyTokenByRPC_Run 期望错误但是没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("TestVerifyTokenByRPC_Run 错误 = %v", err)
				return
			}

			if resp == nil {
				t.Error("TestVerifyTokenByRPC_Run 响应不应该为空")
				return
			}
		})
	}
}
