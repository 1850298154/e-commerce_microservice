package auth

import (
	"bytes"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	// "github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestDeliverTokenByRPC(t *testing.T) {
	h := server.Default()
	h.POST("/auth/token", DeliverTokenByRPC)
	path := "/auth/token"                                     // todo: you can customize query
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := ut.Header{}                                     // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestVerifyTokenByRPC(t *testing.T) {
	h := server.Default()
	h.POST("/auth/renew", VerifyTokenByRPC)
	path := "/auth/renew"                                     // todo: you can customize query
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := ut.Header{}                                     // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestRenewTokenByRPC(t *testing.T) {
	h := server.Default()
	h.POST("/auth/verify", RenewTokenByRPC)
	path := "/auth/verify"                                    // todo: you can customize query
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := ut.Header{}                                     // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "POST", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}
