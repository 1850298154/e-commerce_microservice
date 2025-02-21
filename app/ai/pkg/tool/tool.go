package tool

import (
	"2501YTC/app/ai/infra/rpc"
	rpccart "2501YTC/rpc_gen/kitex_gen/cart"
	rpccheckout "2501YTC/rpc_gen/kitex_gen/checkout"
	rpcorder "2501YTC/rpc_gen/kitex_gen/order"
	rpcpayment "2501YTC/rpc_gen/kitex_gen/payment"
	rpcproduct "2501YTC/rpc_gen/kitex_gen/product"
	"context"
	"encoding/json"
	"github.com/cloudwego/kitex/pkg/klog"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func GetSearchOrdersTool() tool.InvokableTool {
	return &SearchOrdersTool{}
}

type SearchOrdersParam struct {
	UserID uint32 `json:"user_id"`
	//ProductName string `json:"product_name"`
}

type SearchOrdersResult struct {
	OrderId      string                `json:"order_id"`
	UserId       uint32                `json:"user_id"`
	UserCurrency string                `json:"user_currency"`
	Email        string                `json:"email"`
	CreatedAt    time.Time             `json:"created_at"`
	OrderItems   []*rpcorder.OrderItem `json:"order_items"`
	OrderState   string                `json:"orderState"`
}

type SearchOrdersTool struct {
}

func (s *SearchOrdersTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "query_orders",
		Desc: "根据user_id查询订单，返回订单ID、用户ID、用户货币、电子邮件、创建日期、订单商品和订单状态等详细信息。",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"user_id": {
				Type:     "number",
				Desc:     "The id of the user",
				Required: true,
			},
			//"product_name": {
			//	Type:     "string",
			//	Desc:     "The name of the product",
			//	Required: true,
			//},
		}),
	}, nil
}

func (s *SearchOrdersTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &SearchOrdersParam{}
	err := json.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		klog.Errorf("Unmarshal json err: %v", err)
		return "", err
	}

	ordersResp, err := rpc.OrderClient.ListOrder(ctx, &rpcorder.ListOrderReq{
		UserId: p.UserID,
	})
	if err != nil {
		klog.Errorf("ListOrder err: %v", err)
		return "", err
	}
	var orderList []SearchOrdersResult
	for _, resp := range ordersResp.Orders {
		orderList = append(orderList, SearchOrdersResult{
			OrderId:      resp.Order.OrderId,
			UserId:       p.UserID,
			UserCurrency: resp.Order.UserCurrency,
			Email:        resp.Order.Email,
			CreatedAt:    convertInt32ToTime(resp.Order.CreatedAt),
			OrderItems:   resp.Order.OrderItems,
			OrderState:   resp.OrderState,
		})
	}
	resp, err := json.Marshal(orderList)
	if err != nil {
		klog.Errorf("Marshal err: %v", err)
		return "", err
	}
	return string(resp), nil
}

func GetProductTool() tool.InvokableTool {
	return &QueryProductTool{}
}

type QueryProductsParam struct {
	//WantProductName string `json:"want_product_name"`
	ProductId uint32 `json:"product_id"`
}

type QueryProductTool struct{}

func (q *QueryProductTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_products",
		Desc: "根据product_id查询商品详细信息，返回商品名称product_name和其他相关信息。",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"product_id": {
				Type:     "number",
				Desc:     "The id of the product",
				Required: true,
			},
		}),
	}, nil
}

func (q *QueryProductTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &QueryProductsParam{}
	err := json.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		klog.Errorf("Unmarshal json err: %v", err)
		return "", err
	}

	product, err := rpc.ProductClient.GetProduct(ctx, &rpcproduct.GetProductReq{
		Id: p.ProductId,
	})
	if err != nil {
		klog.Errorf("GetProduct err: %v", err)
		return "", err
	}

	resp, err := json.Marshal(product.Product)
	if err != nil {
		klog.Errorf("Marshal err: %v", err)
		return "", err
	}
	return string(resp), nil
}

func GetSearchProductTool() tool.InvokableTool {
	return &SearchProductsTool{}
}

func GetAddToCartTool() tool.InvokableTool {
	return &AddToCartTool{}
}

func GetCheckoutTool() tool.InvokableTool {
	return &CheckoutTool{}
}

type SearchProductParam struct {
	ProductName string `json:"product_name"`
	Quantity    int32  `json:"quantity"`
	Topn        int64  `json:"topn"`
}

type SearchProductsTool struct{}

func (s *SearchProductsTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "search_products",
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

func (s *SearchProductsTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
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
	rests, err := rpc.ProductClient.SearchProductsByName(ctx, &rpcproduct.SearchProductsReq{
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

type AddToCartParam struct {
	UserID    uint32 `json:"user_id"`
	ProductId uint32 `json:"id"`
	Quantity  int32  `json:"quantity"`
}
type AddToCartTool struct{}

func (a *AddToCartTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "add_products_to_cart",
		Desc: "将商品添加到购物车",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"user_id": {
				Type:     "number",
				Desc:     "The id of user",
				Required: true,
			},
			"id": {
				Type:     "number",
				Desc:     "The id of one product",
				Required: true,
			},
			"quantity": {
				Type:     "number",
				Desc:     "the number of products that the user want to buy",
				Required: true,
			},
		}),
	}, nil
}

func (a *AddToCartTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &AddToCartParam{}
	err := json.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		return "", err
	}

	//userId, err := stringToUint32(p.UserID)
	//if err != nil {
	//	klog.Errorf("Unmarshal user id err: %v", err)
	//	return "", err
	//}

	// 请求后端服务
	_, err = rpc.CartClient.AddItem(ctx, &rpccart.AddItemReq{
		Item: &rpccart.CartItem{
			ProductId: p.ProductId,
			Quantity:  p.Quantity,
		},
		UserId: p.UserID,
	})
	if err != nil {
		return "", err
	}

	cart, err := rpc.CartClient.GetCart(ctx, &rpccart.GetCartReq{
		UserId: p.UserID,
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

type CheckoutTool struct{}

func (c *CheckoutTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "checkout",
		Desc: "根据用户购物车的信息，进行结算，创建订单，并返回创建好的订单信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"user_id": {
				Type:     "number",
				Desc:     "The id of user",
				Required: true,
			},
		}),
	}, nil
}

func (c *CheckoutTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &rpcorder.PlaceOrderReq{}
	err := json.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		return "", err
	}

	// 请求后端服务
	checkoutResp, err := rpc.CheckoutClient.Checkout(ctx, &rpccheckout.CheckoutReq{
		UserId:    p.UserId,
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
			CreditCardNumber:          "5302079249905900",
			CreditCardCvv:             123,
			CreditCardExpirationMonth: 12,
			CreditCardExpirationYear:  2025,
		},
	})
	if err != nil {
		return "", err
	}

	orderResp, err := rpc.OrderClient.GetOrder(ctx, &rpcorder.GetOrderReq{
		UserId:  p.UserId,
		OrderId: checkoutResp.OrderId,
	})
	if err != nil {
		klog.Error(err)
		return "", err
	}
	// 序列化结果
	res, err := json.Marshal(orderResp.Order.Order)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func convertInt32ToTime(timestamp int32) time.Time {
	// 将 int32 转换为 int64，因为 time.Unix 需要 int64 类型的参数
	seconds := int64(timestamp)
	// 使用 time.Unix 创建 time.Time 对象
	return time.Unix(seconds, 0)
}
