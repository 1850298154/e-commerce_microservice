package auth

import (
	"2501YTC/app/auth/biz/dal"
	"2501YTC/app/gateway/hertz_gen/gateway/auth"
	"2501YTC/app/gateway/infra/rpc"
	"bytes"
	"testing"

	"github.com/goccy/go-json"
	"github.com/joho/godotenv"

	"github.com/cloudwego/hertz/pkg/app/server"
	// "github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestDeliverTokenByRPC(t *testing.T) {
	_ = godotenv.Load("../../../.env")
	rpc.InitClient()
	h := server.Default()
	h.POST("/auth/token", DeliverTokenByRPC)
	// if rpc.AuthClient == nil {
	//	log.Fatal("rpc.AuthClient is nil, please check initialization")
	// }
	path := "/auth/token" // todo: you can customize query
	reqBody := auth.DeliverTokenReq{
		UserId: 1434424,
	}

	jsonBody, _ := json.Marshal(reqBody)
	body := &ut.Body{Body: bytes.NewBuffer(jsonBody), Len: len(jsonBody)} // todo: you can customize body
	header := ut.Header{
		Key:   "Content-Type",
		Value: "application/json",
	}
	// todo: you can customize header
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestVerifyTokenByRPC(t *testing.T) {
	_ = godotenv.Load("../../../.env")
	dal.Init()
	rpc.InitClient()
	h := server.Default()
	h.POST("/auth/verify", VerifyTokenByRPC)
	path := "/auth/verify"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjAsIlJvbGUiOjEsImV4cCI6MTczOTM1NDc3OSwianRpIjoiMDRiNGExN2ItOWY0OS00YjA5LWIxZjAtZGE4YjA5OWEzZTgzIiwiaWF0IjoxNzM5MzUxMTc5LCJpc3MiOiJnb21hbGwifQ._6kU5z5qOCnyVAw2Jk5U_Ki381vB4gnqCQ5t_6n73dg"
	refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjAsIlJvbGUiOjEsImV4cCI6MTczOTk1NTk3OSwianRpIjoiMDRiNGExN2ItOWY0OS00YjA5LWIxZjAtZGE4YjA5OWEzZTgzIiwiaWF0IjoxNzM5MzUxMTc5LCJpc3MiOiJnb21hbGwifQ.G-mp8wm90uGbph899GT2o4bCwkRCNCiVlEce8fnxQZ8"
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := []ut.Header{
		{
			Key:   "Authorization",
			Value: "Bearer " + token,
		},
		{
			Key:   "X-Refresh-Token",
			Value: "Bearer " + refreshToken,
		},
	} // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "POST", path, body, header...)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestRenewTokenByRPC(t *testing.T) {
	_ = godotenv.Load("../../../.env")
	dal.Init()
	rpc.InitClient()
	h := server.Default()
	h.POST("/auth/renew", RenewTokenByRPC)
	path := "/auth/renew"
	refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEyMzQ1LCJSb2xlIjoxLCJleHAiOjE3Mzk5ODMxNjcsImp0aSI6IjhjNWJlYzYzLWNiZjEtNGYyMC04YWM4LTg4MTcyNTQ2MjdkZiIsImlhdCI6MTczOTM3Nzk4OSwiaXNzIjoiZ29tYWxsIn0.K_y3fXtNJ3Ccjg9VtbkAWX-vpTlPLZizIUWWXitTBMg" // todo: you can customize query
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1}                                                                                                                                                                                                          // todo: you can customize body
	header := ut.Header{
		Key:   "X-Refresh-Token",
		Value: "Bearer " + refreshToken,
	} // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}
