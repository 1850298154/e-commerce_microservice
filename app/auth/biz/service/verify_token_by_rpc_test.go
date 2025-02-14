package service

import (
	"2501YTC/app/auth/biz/dal/mysql"
	"2501YTC/app/auth/biz/dal/redis"
	"context"
	"testing"

	"github.com/joho/godotenv"

	auth "2501YTC/rpc_gen/kitex_gen/auth"
)

func TestVerifyTokenByRPC_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	mysql.Init()
	redis.Init()
	ctx := context.Background()
	s := NewVerifyTokenByRPCService(ctx)
	// init req and assert value
	// token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjQsIlJvbGUiOjQsImV4cCI6MTczOTI5NjI5MiwianRpIjoiYzlkYTdjNzgtNjM1NC00Njc4LTg1ZDgtOWM3YzRlYTYwZTQ5IiwiaWF0IjoxNzM5MjkyNjkyLCJpc3MiOiJnb21hbGwifQ.cGXL0DX9A927-H-cVh2DUdOnAVcvUUcFSUiI_0-9608"
	// refresh_token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjQsIlJvbGUiOjQsImV4cCI6MTczOTg5NzQ5MiwianRpIjoiYzlkYTdjNzgtNjM1NC00Njc4LTg1ZDgtOWM3YzRlYTYwZTQ5IiwiaWF0IjoxNzM5MjkyNjkyLCJpc3MiOiJnb21hbGwifQ.Xb55VrUlnoguVqgbIB2Ey_xynMDEpwgSi46SHf_1MPo"
	req := &auth.VerifyTokenReq{
		Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjQsIlJvbGUiOjQsImV4cCI6MTczOTI5NjI5MiwianRpIjoiYzlkYTdjNzgtNjM1NC00Njc4LTg1ZDgtOWM3YzRlYTYwZTQ5IiwiaWF0IjoxNzM5MjkyNjkyLCJpc3MiOiJnb21hbGwifQ.cGXL0DX9A927-H-cVh2DUdOnAVcvUUcFSUiI_0-9608",
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjQsIlJvbGUiOjQsImV4cCI6MTczOTg5NzQ5MiwianRpIjoiYzlkYTdjNzgtNjM1NC00Njc4LTg1ZDgtOWM3YzRlYTYwZTQ5IiwiaWF0IjoxNzM5MjkyNjkyLCJpc3MiOiJnb21hbGwifQ.Xb55VrUlnoguVqgbIB2Ey_xynMDEpwgSi46SHf_1MPo",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
	// todo: edit your unit test
}
