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
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjI2ODM2NjUzMzcsIlJvbGUiOjIsImlzcyI6ImdvbWFsbCIsImV4cCI6MTc0MDEyNTc3MywiaWF0IjoxNzQwMTIyMTczLCJqdGkiOiIyYThmNTFmMC1mMTkzLTQ3YmEtOTRiMC00ZTlmYWZlMDE0YmIifQ.q0RDqifqQwoWq3NlzWXRK_9WOkDbda940d0AoBqHT7s"
	refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjI2ODM2NjUzMzcsIlJvbGUiOjIsImlzcyI6ImdvbWFsbCIsImV4cCI6MTc0MDcyNjk3MywiaWF0IjoxNzQwMTIyMTczLCJqdGkiOiIyYThmNTFmMC1mMTkzLTQ3YmEtOTRiMC00ZTlmYWZlMDE0YmIifQ.BPNBsWzBX84lkL1V9J9DDLVfUL3lJZMV8MWT1eOQ0JQ"
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := []ut.Header{
		{
			Key:   "Authorization",
			Value: token,
		},
		{
			Key:   "X-Refresh-Token",
			Value: refreshToken,
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
