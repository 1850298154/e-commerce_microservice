package order

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"2501YTC/app/gateway/infra/rpc"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder(t *testing.T) {
	rpc.InitClient()
	h := server.Default()
	h.POST("/orders", PlaceOrder)

	testCases := []struct {
		name       string
		reqBody    string
		wantStatus int
	}{
		{
			name: "valid order",
			reqBody: `{
				"user_id": 124,
				"user_currency": "USD",
				"address": {
					"street_address": "123 Main St",
					"city": "Boston", 
					"state": "MA",
					"country": "USA",
					"zip_code": 12345
				},
				"email": "test@example.com",
				"order_items": [
					{
						"product_id": 1,
						"quantity": 2,
						"cost": 9.99
					}
				]
			}`,
			wantStatus: 200,
		},
		{
			name:       "empty body",
			reqBody:    "",
			wantStatus: 400,
		},
		{
			name: "invalid order - missing required fields",
			reqBody: `{
				"user_id": 123
			}`,
			wantStatus: 400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := &ut.Body{
				Body: bytes.NewBufferString(tc.reqBody),
				Len:  len(tc.reqBody),
			}
			header := ut.Header{
				Key:   "Content-Type",
				Value: "application/json",
			}

			w := ut.PerformRequest(h.Engine, "POST", "/orders", body, header)
			resp := w.Result()

			assert.Equal(t, tc.wantStatus, resp.StatusCode())
			t.Logf("Response body: %s", string(resp.Body()))
		})
	}
}

func TestOrderLifecycle(t *testing.T) {
	rpc.InitClient()
	h := server.Default()

	// 注册所有路由
	h.POST("/orders", PlaceOrder)
	h.GET("/orders", ListOrder)
	h.PUT("/orders/:order_id/paid", MarkOrderPaid)
	h.PUT("/orders/:order_id", UpdateOrder)
	h.DELETE("/orders/:order_id", CancelOrder)

	// 1. 创建订单
	placeOrderBody := `{
        "user_id": 124,
        "user_currency": "USD",
        "address": {
            "street_address": "123 Main St",
            "city": "Boston",
            "state": "MA",
            "country": "USA",
            "zip_code": 12345
        },
        "email": "test@example.com",
        "order_items": [
            {
                "product_id": 1,
                "quantity": 2,
                "cost": 9.99
            }
        ]
    }`

	w := ut.PerformRequest(h.Engine, "POST", "/orders",
		&ut.Body{
			Body: bytes.NewBufferString(placeOrderBody),
			Len:  len(placeOrderBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp := w.Result()
	assert.Equal(t, RPCCallSuccess, resp.StatusCode())
	t.Logf("Place order response: %s", string(resp.Body()))

	// 从响应中获取order_id
	var placeOrderResp struct {
		OrderId string `json:"order"`
	}
	if err := json.Unmarshal(resp.Body(), &placeOrderResp); err != nil {
		t.Fatalf("Failed to unmarshal place order response: %v", err)
	}
	orderId := placeOrderResp.OrderId

	// 2. 列出订单
	w = ut.PerformRequest(h.Engine, "GET", "/orders?user_id=124", nil)
	resp = w.Result()
	assert.Equal(t, RPCCallSuccess, resp.StatusCode())
	t.Logf("List orders response: %s", string(resp.Body()))

	// 3. 标记订单已支付
	markPaidBody := `{"user_id": 124}`
	w = ut.PerformRequest(h.Engine, "PUT",
		fmt.Sprintf("/orders/%s/paid", orderId),
		&ut.Body{
			Body: bytes.NewBufferString(markPaidBody),
			Len:  len(markPaidBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp = w.Result()
	assert.Equal(t, RPCCallSuccess, resp.StatusCode())
	t.Logf("Mark order paid response: %s", string(resp.Body()))

	// 4. 更新订单
	updateOrderBody := `{
        "user_id": 124,
        "new_address": {
            "street_address": "456 New St",
            "city": "New York",
            "state": "NY",
            "country": "USA",
            "zip_code": 54321
        },
        "new_email": "newemail@example.com",
        "new_order_items": [
            {
                "product_id": 1,
                "quantity": 3,
                "cost": 9.99
            }
        ]
    }`
	w = ut.PerformRequest(h.Engine, "PUT",
		fmt.Sprintf("/orders/%s", orderId),
		&ut.Body{
			Body: bytes.NewBufferString(updateOrderBody),
			Len:  len(updateOrderBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp = w.Result()
	assert.Equal(t, RPCCallSuccess, resp.StatusCode())
	t.Logf("Update order response: %s", string(resp.Body()))

	// 5. 取消订单
	cancelOrderBody := `{
        "user_id": 124,
        "timed_cancel": false,
        "cancel_time": 0
    }`
	w = ut.PerformRequest(h.Engine, "DELETE",
		fmt.Sprintf("/orders/%s", orderId),
		&ut.Body{
			Body: bytes.NewBufferString(cancelOrderBody),
			Len:  len(cancelOrderBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp = w.Result()
	assert.Equal(t, RPCCallSuccess, resp.StatusCode())
	t.Logf("Cancel order response: %s", string(resp.Body()))
}

func TestListOrder(t *testing.T) {
	h := server.Default()
	h.GET("/orders", ListOrder)
	path := "/orders"                                         // todo: you can customize query
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := ut.Header{}                                     // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "GET", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestMarkOrderPaid(t *testing.T) {
	h := server.Default()
	h.PUT("/orders/:order_id/paid", MarkOrderPaid)
	path := "/orders/:order_id/paid"                          // todo: you can customize query
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := ut.Header{}                                     // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "PUT", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestUpdateOrder(t *testing.T) {
	h := server.Default()
	h.PUT("/orders/:order_id", UpdateOrder)
	path := "/orders/:order_id"                               // todo: you can customize query
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := ut.Header{}                                     // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "PUT", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}

func TestCancelOrder(t *testing.T) {
	h := server.Default()
	h.DELETE("/orders/:order_id", CancelOrder)
	path := "/orders/:order_id"                               // todo: you can customize query
	body := &ut.Body{Body: bytes.NewBufferString(""), Len: 1} // todo: you can customize body
	header := ut.Header{}                                     // todo: you can customize header
	w := ut.PerformRequest(h.Engine, "DELETE", path, body, header)
	resp := w.Result()
	t.Log(string(resp.Body()))

	// todo edit your unit test.
	// assert.DeepEqual(t, 200, resp.StatusCode())
	// assert.DeepEqual(t, "null", string(resp.Body()))
}
