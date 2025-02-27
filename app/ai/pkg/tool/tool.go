package tool

import (
	"context"
	"encoding/json"
	"time"

	"2501YTC/app/ai/infra/rpc"
	rpccart "2501YTC/rpc_gen/kitex_gen/cart"
	rpccheckout "2501YTC/rpc_gen/kitex_gen/checkout"
	rpcorder "2501YTC/rpc_gen/kitex_gen/order"
	rpcpayment "2501YTC/rpc_gen/kitex_gen/payment"
	rpcproduct "2501YTC/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func GetSearchOrdersTool() tool.InvokableTool {
	return &SearchOrdersTool{}
}

type SearchOrdersParam struct {
	UserID uint32 `json:"user_id"`
}

type OrderItem struct {
	ProductId   uint32  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int32   `json:"quantity"`
	Cost        float32 `json:"cost"`
}

type SearchOrdersResult struct {
	OrderId      string      `json:"order_id"`
	UserId       uint32      `json:"user_id"`
	UserCurrency string      `json:"user_currency"`
	Email        string      `json:"email"`
	CreatedAt    time.Time   `json:"created_at"`
	OrderItems   []OrderItem `json:"order_items"`
	OrderState   string      `json:"orderState"`
}

type SearchOrdersTool struct{}

func (s *SearchOrdersTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "query_orders",
		Desc: "Query orders based on user_id and return detailed information such as orderID, userID, userCurrency, email, createAt, orderItems, and orderState",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"user_id": {
				Type:     "number",
				Desc:     "The id of the user",
				Required: true,
			},
		}),
	}, nil
}

func (s *SearchOrdersTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &SearchOrdersParam{}
	err := json.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		klog.Error(err)
		return "", err
	}

	// 查询用户的所有订单信息
	ordersResp, err := rpc.OrderClient.ListOrder(ctx, &rpcorder.ListOrderReq{
		UserId: p.UserID,
	})
	if err != nil {
		klog.Errorf("ListOrder err: %v", err)
		return "", err
	}

	// 将结果组合成SearchOrdersResult结构体
	orderList := make([]SearchOrdersResult, 0)
	for _, resp := range ordersResp.Orders {
		orderItems := make([]OrderItem, 0)
		for _, item := range resp.Order.OrderItems {
			product, err := rpc.ProductClient.GetProduct(ctx, &rpcproduct.GetProductReq{Id: item.Item.ProductId})
			if err != nil {
				klog.Errorf("GetProduct err: %v", err)
				return "", err
			}
			orderItems = append(orderItems, OrderItem{
				ProductId:   item.Item.ProductId,
				ProductName: product.Product.Name,
				Quantity:    item.Item.Quantity,
				Cost:        item.Cost,
			})
		}
		orderList = append(orderList, SearchOrdersResult{
			OrderId:      resp.Order.OrderId,
			UserId:       p.UserID,
			UserCurrency: resp.Order.UserCurrency,
			Email:        resp.Order.Email,
			CreatedAt:    convertInt32ToTime(resp.Order.CreatedAt),
			OrderItems:   orderItems,
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
		Name: "search products",
		Desc: "query the specified product based on the product name",
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

	// 调用商品服务查找特定名称的商品
	rests, err := rpc.ProductClient.SearchProductsByName(ctx, &rpcproduct.SearchProductsByNameReq{
		Query:    p.ProductName,
		Page:     1,
		PageSize: p.Topn,
		Flag:     false,
	})
	if err != nil {
		klog.Errorf("SearchProductsByName err: %v", err)
		return "", err
	}

	// 序列化结果
	res, err := json.Marshal(rests.Results)
	if err != nil {
		klog.Error(err)
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
		Name: "add products to cart",
		Desc: "add the selected items to the shopping cart.",
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
		klog.Error(err)
		return "", err
	}

	// 调用购物车服务将商品添加到购物车
	_, err = rpc.CartClient.AddItem(ctx, &rpccart.AddItemReq{
		Item: &rpccart.CartItem{
			ProductId: p.ProductId,
			Quantity:  p.Quantity,
		},
		UserId: p.UserID,
	})
	if err != nil {
		klog.Errorf("AddItem err: %v", err)
		return "", err
	}

	// 返回购物车信息
	cart, err := rpc.CartClient.GetCart(ctx, &rpccart.GetCartReq{
		UserId: p.UserID,
	})
	if err != nil {
		klog.Errorf("GetCart err: %v", err)
		return "", err
	}
	// 序列化结果
	res, err := json.Marshal(cart.Cart.Items)
	if err != nil {
		klog.Error(err)
		return "", err
	}

	return string(res), nil
}

type CheckoutTool struct{}

func (c *CheckoutTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "checkout",
		Desc: "settle the payment based on the items in the user's shopping cart, create an order, and return the created order information.",
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

	// 调用结算服务进行订单结算
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
		klog.Errorf("checkout failed: %s", err)
		return "", err
	}

	// 获取下单后的订单信息
	orderResp, err := rpc.OrderClient.GetOrder(ctx, &rpcorder.GetOrderReq{
		UserId:  p.UserId,
		OrderId: checkoutResp.OrderId,
	})
	if err != nil {
		klog.Error(err)
		return "", err
	}

	// 将订单信息组合成SearchOrdersResult结构体返回
	orderItems := make([]OrderItem, 0)
	for _, item := range orderResp.Order.Order.OrderItems {
		product, err := rpc.ProductClient.GetProduct(ctx, &rpcproduct.GetProductReq{Id: item.Item.ProductId})
		if err != nil {
			klog.Error("get product name failed: %s", err)
			return "", err
		}
		orderItems = append(orderItems, OrderItem{
			ProductId:   item.Item.ProductId,
			ProductName: product.Product.Name,
			Quantity:    item.Item.Quantity,
			Cost:        item.Cost,
		})
	}
	order := SearchOrdersResult{
		OrderId:      orderResp.Order.Order.OrderId,
		UserId:       orderResp.Order.Order.UserId,
		UserCurrency: orderResp.Order.Order.UserCurrency,
		Email:        orderResp.Order.Order.Email,
		CreatedAt:    convertInt32ToTime(orderResp.Order.Order.CreatedAt),
		OrderItems:   orderItems,
		OrderState:   orderResp.Order.OrderState,
	}

	res, err := json.Marshal(order)
	if err != nil {
		klog.Error(err)
		return "", err
	}

	return string(res), nil
}

func convertInt32ToTime(timestamp int32) time.Time {
	seconds := int64(timestamp)

	return time.Unix(seconds, 0)
}
