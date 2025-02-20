package tool

import (
	"2501YTC/app/ai/infra/rpc"
	"2501YTC/app/gateway/hertz_gen/cart"
	rpccart "2501YTC/rpc_gen/kitex_gen/cart"
	rpcorder "2501YTC/rpc_gen/kitex_gen/order"
	rpcproduct "2501YTC/rpc_gen/kitex_gen/product"
	rpccheckout "2501YTC/rpc_gen/kitex_gen/checkout"
	rpcpayment "2501YTC/rpc_gen/kitex_gen/payment"
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type SearchProductParam struct {
	ProductName string `json:"product_name"`
	Quantity    int32  `json:"quantity"`
	Topn        int64  `json:"topn"`
}

type ToolSearchProducts struct{}

func (t *ToolSearchProducts) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "search products",
		Desc: "查询指定商品",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"product_name": {
				Type:     "string",
				Desc:     "The name of one product",
				Required: true,
			},
			"topn": {
				Type: "number",
				Desc: "top n products sorted by prices",
			},
		}),
	}, nil
}

func (t *ToolSearchProducts) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &SearchProductParam{}
	err := json.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		return "", err
	}

	if p.Topn == 0 {
		p.Topn = 1
	}

	// 请求后端服务
	rests, err := rpc.ProductClient.SearchProducts(ctx, &rpcproduct.SearchProductsReq{
		Query:    p.ProductName,
		Page:     1,
		PageSize: p.Topn,
	})
	if err != nil {
		return "", err
	}

	// 序列化结果
	res, err := json.Marshal(rests.Results)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

//type AddToCartParam struct {
//	ProductId string `json:"product_id"`
//	Quantity  int32  `json:"quantity"`
//}
type ToolAddToCart struct{}

func (t *ToolAddToCart) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "add products to cart",
		Desc: "将商品添加到购物车",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"product_id": {
				Type:     "string",
				Desc:     "The id of one product",
				Required: true,
			},
		}),
	}, nil
}

func (t *ToolAddToCart) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &rpccart.CartItem{}
	err := json.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		return "", err
	}

	userId := 1

	// 请求后端服务
	_, err = rpc.CartClient.AddItem(ctx, &rpccart.AddItemReq{
		Item:   p,
		UserId: uint32(userId),
	})
	if err != nil {
		return "", err
	}

	cart, err := rpc.CartClient.GetCart(ctx, &rpccart.GetCartReq{
		UserId: uint32(userId),
	})
	if err != nil {
		return "", err
	}
	// 序列化结果
	res, err := json.Marshal(cart.Cart.Items)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

type ToolCheckout struct{}

func (t *ToolCheckout) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "create order",
		Desc: "根据用户购物车的信息，创建订单",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"product_id": {
				Type:     "string",
				Desc:     "The id of one product",
				Required: true,
			},
		}),
	}, nil
}

func (t *ToolCheckout) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &rpcorder.PlaceOrderReq{}
	err := json.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		return "", err
	}

	userId := 1

	// 请求后端服务
	checkoutResp, err := rpc.CheckoutClient.Checkout(ctx, &rpccheckout.CheckoutReq{
		UserId:    uint32(userId),
		Firstname: "user",
		Lastname:  "user",
		Address: &rpccheckout.Address{
			StreetAddress: "123 Main St",
			City:          "Beijing",
			State:         "Beijing",
			Country:       "China",
			ZipCode:       "0",
		},
		Email: "user@example.com",
		CreditCard: &rpcpayment.CreditCardInfo{
			CreditCardNumber:          "",
			CreditCardCvv:             1,
			CreditCardExpirationMonth: 2,
			CreditCardExpirationYear:  3,
		},
	})
	if err != nil {
		return "", err
	}

	orderResp, err := rpc.OrderClient.ListOrder(ctx, &rpcorder.ListOrderReq{
		UserId: uint32(userId),
	})
	// 序列化结果
	res, err := json.Marshal()
	if err != nil {
		return "", err
	}

	return string(res), nil
}
