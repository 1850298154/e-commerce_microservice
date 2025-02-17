package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"2501YTC/app/gateway/infra/rpc"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder(t *testing.T) {
	rpc.InitClient()
	h := server.Default()
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Set("user_id", uint32(124))
	})

	h.POST("/orders", PlaceOrder)

	testCases := []struct {
		name       string
		reqBody    string
		wantStatus int
	}{
		{
			name: "valid order",
			reqBody: `{
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

func TestListOrder(t *testing.T) {
	rpc.InitClient()
	h := server.Default()

	// 模拟用户登录中间件
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Set("user_id", uint32(124))
	})

	// 注册需要的路由
	h.POST("/orders", PlaceOrder)
	h.GET("/orders", ListOrder)

	// 1. 首先创建一个测试订单
	createOrderBody := `{
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

	// 创建订单
	w := ut.PerformRequest(h.Engine, "POST", "/orders",
		&ut.Body{
			Body: bytes.NewBufferString(createOrderBody),
			Len:  len(createOrderBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp := w.Result()
	assert.Equal(t, RPCCallSuccess, resp.StatusCode())

	// 2. 测试正常列出订单
	w = ut.PerformRequest(h.Engine, "GET", "/orders",
		&ut.Body{Body: bytes.NewBufferString(""), Len: 0},
		ut.Header{},
	)
	resp = w.Result()
	assert.Equal(t, RPCCallSuccess, resp.StatusCode())
	t.Logf("List orders response: %s", string(resp.Body()))

	// 验证响应中包含订单数据
	var listResp struct {
		Orders []any `json:"orders"`
	}
	err := json.Unmarshal(resp.Body(), &listResp)
	assert.Nil(t, err)
	assert.Greater(t, len(listResp.Orders), 0)

	// 3. 测试列出不存在用户的订单
	h2 := server.Default()
	h2.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Set("user_id", uint32(999999)) // 使用一个不存在的用户ID
	})
	h2.GET("/orders", ListOrder)

	w = ut.PerformRequest(h2.Engine, "GET", "/orders",
		&ut.Body{Body: bytes.NewBufferString(""), Len: 0},
		ut.Header{},
	)
	resp = w.Result()
	assert.Equal(t, RPCCallSuccess, resp.StatusCode())

	var emptyListResp struct {
		Orders []any `json:"orders"`
	}
	err = json.Unmarshal(resp.Body(), &emptyListResp)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(emptyListResp.Orders))
}

func TestMarkOrderPaid(t *testing.T) {
	rpc.InitClient()
	h := server.Default()

	// 模拟用户登录
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Set("user_id", uint32(124))
	})

	// 注册需要的路由
	h.POST("/orders", PlaceOrder)
	h.PUT("/orders/:order_id/paid", MarkOrderPaid)

	// 1. 首先创建一个测试订单
	createOrderBody := `{
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
			Body: bytes.NewBufferString(createOrderBody),
			Len:  len(createOrderBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp := w.Result()
	assert.Equal(t, RPCCallSuccess, resp.StatusCode())

	// 解析订单ID
	var placeOrderResp struct {
		OrderId string `json:"order"`
	}
	err := json.Unmarshal(resp.Body(), &placeOrderResp)
	assert.Nil(t, err)
	orderId := placeOrderResp.OrderId

	// 2. 测试标记其他用户的订单
	h2 := server.Default()
	h2.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Set("user_id", uint32(999)) // 使用不同的用户ID
	})
	h2.PUT("/orders/:order_id/paid", MarkOrderPaid)

	w = ut.PerformRequest(h2.Engine, "PUT",
		fmt.Sprintf("/orders/%s/paid", orderId),
		&ut.Body{
			Body: bytes.NewBufferString(`{"user_id": 999}`),
			Len:  len(`{"user_id": 999}`),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp = w.Result()
	assert.Equal(t, 500, resp.StatusCode())

	// 3. 测试正常标记订单已支付
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

	// 4. 测试标记不存在的订单
	w = ut.PerformRequest(h.Engine, "PUT",
		"/orders/non-existent-order/paid",
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
	assert.Equal(t, 500, resp.StatusCode())

	// 5. 测试重复标记已支付的订单
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
	assert.Equal(t, 500, resp.StatusCode())
}

func TestUpdateOrder(t *testing.T) {
	rpc.InitClient()
	h := server.Default()

	// 模拟用户登录
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Set("user_id", uint32(124))
	})

	// 注册需要的路由
	h.POST("/orders", PlaceOrder)
	h.PUT("/orders/:order_id", UpdateOrder)
	h.DELETE("/orders/:order_id", CancelOrder)

	// 1. 首先创建一个测试订单
	createOrderBody := `{
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
			Body: bytes.NewBufferString(createOrderBody),
			Len:  len(createOrderBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp := w.Result()
	assert.Equal(t, RPCCallSuccess, resp.StatusCode())

	// 解析订单ID
	var placeOrderResp struct {
		OrderId string `json:"order"`
	}
	err := json.Unmarshal(resp.Body(), &placeOrderResp)
	assert.Nil(t, err)
	orderId := placeOrderResp.OrderId

	// 2. 测试正常更新订单
	updateOrderBody := `{
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

	// 3. 测试更新不存在的订单
	w = ut.PerformRequest(h.Engine, "PUT",
		"/orders/non-existent-order",
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
	assert.Equal(t, 500, resp.StatusCode())
	t.Logf("Update order response: %s", string(resp.Body()))

	// 4. 测试更新其他用户的订单
	h2 := server.Default()
	h2.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Set("user_id", uint32(999)) // 使用不同的用户ID
	})
	h2.PUT("/orders/:order_id", UpdateOrder)

	wrongUserBody := `{
        "new_email": "wrong@example.com"
    }`
	w = ut.PerformRequest(h2.Engine, "PUT",
		fmt.Sprintf("/orders/%s", orderId),
		&ut.Body{
			Body: bytes.NewBufferString(wrongUserBody),
			Len:  len(wrongUserBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp = w.Result()
	assert.Equal(t, 500, resp.StatusCode())
	t.Logf("Update order response: %s", string(resp.Body()))

	// 5. 测试空更新
	emptyUpdateBody := `{
    }`
	w = ut.PerformRequest(h.Engine, "PUT",
		fmt.Sprintf("/orders/%s", orderId),
		&ut.Body{
			Body: bytes.NewBufferString(emptyUpdateBody),
			Len:  len(emptyUpdateBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp = w.Result()
	assert.Equal(t, 500, resp.StatusCode())
	t.Logf("Update order response: %s", string(resp.Body()))

	// 6. 测试更新已取消的订单
	// 首先取消订单
	cancelOrderBody := `{
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

	// 尝试更新已取消的订单
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
	assert.Equal(t, 500, resp.StatusCode())
	t.Logf("Update order response: %s", string(resp.Body()))
}

func TestCancelOrder(t *testing.T) {
	rpc.InitClient()
	h := server.Default()

	// 模拟用户登录
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Set("user_id", uint32(124))
	})

	// 注册需要的路由
	h.POST("/orders", PlaceOrder)
	h.DELETE("/orders/:order_id", CancelOrder)

	// 1. 首先创建一个测试订单
	createOrderBody := `{
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
			Body: bytes.NewBufferString(createOrderBody),
			Len:  len(createOrderBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp := w.Result()
	assert.Equal(t, RPCCallSuccess, resp.StatusCode())

	// 解析订单ID
	var placeOrderResp struct {
		OrderId string `json:"order"`
	}
	err := json.Unmarshal(resp.Body(), &placeOrderResp)
	assert.Nil(t, err)
	orderId := placeOrderResp.OrderId

	// 2. 测试取消其他用户的订单
	h2 := server.Default()
	h2.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Set("user_id", uint32(999)) // 使用不同的用户ID
	})
	h2.DELETE("/orders/:order_id", CancelOrder)

	wrongUserBody := `{
			"timed_cancel": false,
			"cancel_time": 0
		}`
	w = ut.PerformRequest(h2.Engine, "DELETE",
		fmt.Sprintf("/orders/%s", orderId),
		&ut.Body{
			Body: bytes.NewBufferString(wrongUserBody),
			Len:  len(wrongUserBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp = w.Result()
	assert.Equal(t, 500, resp.StatusCode())

	// 3. 测试正常取消订单
	cancelOrderBody := `{
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

	// 4. 测试取消不存在的订单
	w = ut.PerformRequest(h.Engine, "DELETE",
		"/orders/non-existent-order",
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
	assert.Equal(t, 500, resp.StatusCode())

	// 5. 测试重复取消订单
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
	assert.Equal(t, 500, resp.StatusCode())

	// 6. 测试定时取消订单但未提供取消时间
	invalidTimedCancelBody := `{
        "timed_cancel": true,
        "cancel_time": 0
    }`
	w = ut.PerformRequest(h.Engine, "DELETE",
		fmt.Sprintf("/orders/%s", orderId),
		&ut.Body{
			Body: bytes.NewBufferString(invalidTimedCancelBody),
			Len:  len(invalidTimedCancelBody),
		},
		ut.Header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	resp = w.Result()
	assert.Equal(t, 500, resp.StatusCode())
}

func TestOrderLifecycle(t *testing.T) {
	rpc.InitClient()
	h := server.Default()
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Set("user_id", uint32(124))
	})

	// 注册所有路由
	h.POST("/orders", PlaceOrder)
	h.GET("/orders", ListOrder)
	h.PUT("/orders/:order_id/paid", MarkOrderPaid)
	h.PUT("/orders/:order_id", UpdateOrder)
	h.DELETE("/orders/:order_id", CancelOrder)

	// 1. 创建订单
	placeOrderBody := `{
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
	markPaidBody := `{}`
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
