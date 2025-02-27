package ai

import (
	"bytes"
	"encoding/json"
	"testing"

	"2501YTC/app/gateway/infra/rpc"

	"2501YTC/app/gateway/hertz_gen/gateway/ai"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestQueryOrder(t *testing.T) {
	rpc.InitClient()
	h := server.Default()
	h.POST("/ai/query", QueryOrder)
	path := "/ai/query"
	reqBody := ai.QueryOrderReq{
		Content: "查找前天的订单",
	}
	jsonBody, _ := json.Marshal(reqBody)
	body := &ut.Body{Body: bytes.NewBuffer(jsonBody), Len: len(jsonBody)}
	header := ut.Header{
		Key:   "Content-Type",
		Value: "application/json",
	}
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestAutoOrder(t *testing.T) {
	rpc.InitClient()
	h := server.Default()
	h.POST("/ai/place", AutoOrder)
	path := "/ai/place"
	reqBody := ai.PlaceOrderReq{
		Content: "购买2件衬衫和2个小米手环9。",
	}
	jsonBody, _ := json.Marshal(reqBody)
	body := &ut.Body{Body: bytes.NewBufferString(string(jsonBody)), Len: len(jsonBody)}
	header := ut.Header{
		Key:   "Content-Type",
		Value: "application/json",
	}
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}
