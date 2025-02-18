package tool

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"2501YTC/app/ai/infra/rpc"
	rpccart "2501YTC/rpc_gen/kitex_gen/cart"
	rpcorder "2501YTC/rpc_gen/kitex_gen/order"
	rpcproduct "2501YTC/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/samber/lo"
)

type productItem struct {
	ProductName string `json:"product_name"`
	Quantity    int32  `json:"quantity"`
}

func autoOrderFunc(ctx context.Context, param map[string][]productItem) (string, error) {
	productItems := param["products"]
	var cartItems []*rpccart.CartItem = make([]*rpccart.CartItem, 0, len(productItems))
	var orderItems []*rpcorder.OrderItem = make([]*rpcorder.OrderItem, 0, len(productItems))
	for _, item := range productItems {
		resp, err := rpc.ProductClient.SearchProducts(ctx, &rpcproduct.SearchProductsReq{
			Query:    item.ProductName,
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			klog.Errorf("search product %s error: %v", item.ProductName, err)
			return "", err
		}
		sort.Slice(resp.Results, func(i, j int) bool {
			return resp.Results[i].Price < resp.Results[j].Price
		})
		product := resp.Results[0]
		cartItem := &rpccart.CartItem{
			ProductId: product.Id,
			Quantity:  item.Quantity,
		}
		cartItems = append(cartItems, cartItem)
		orderItem := &rpcorder.OrderItem{
			Item: cartItem,
			Cost: product.Price * float32(cartItem.Quantity),
		}
		orderItems = append(orderItems, orderItem)
	}
	userId := 1
	for _, item := range cartItems {
		_, err := rpc.CartClient.AddItem(ctx, &rpccart.AddItemReq{
			UserId: uint32(userId),
			Item:   item,
		})
		if err != nil {
			return "", err
		}
	}
	resp, err := rpc.OrderClient.PlaceOrder(ctx, &rpcorder.PlaceOrderReq{
		UserId:     uint32(userId),
		OrderItems: orderItems,
	})
	if err != nil {
		klog.Errorf("place order error: %v", err)
		return "", err
	}
	return resp.Order.OrderId, nil
}

func autoOrderTool() tool.InvokableTool {
	info := &schema.ToolInfo{
		Name: "自动下单",
		Desc: "根据用户想要购买的商品名称，自动搜索商品，加到购物车并创建相应的订单，返回给用户创建好的订单信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"products": {
				Desc: "用户想要购买的商品数组，包含了商品的名称和数量，例如：[{product_name: 手机, quantity: 1},{product_name: 电脑, quantity: 2}]",
				Type: schema.Array,
				ElemInfo: &schema.ParameterInfo{
					Type: schema.Object,
				},
				Required: true,
			},
		}),
	}

	return utils.NewTool(info, autoOrderFunc)
}

func NewAutoOrderTools() []tool.BaseTool {
	// 初始化 tool
	tools := []tool.BaseTool{
		autoOrderTool(),
	}
	return tools
}

func searchOrdersFunc(ctx context.Context, param map[string][]string) (string, error) {
	var productsName []string
	for _, item := range param["names"] {
		productsName = append(productsName, item)
	}
	userId := 1
	listOrder, err := rpc.OrderClient.ListOrder(ctx, &rpcorder.ListOrderReq{UserId: uint32(userId)})
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(listOrder.Orders)
	if err != nil {
		klog.Errorf("Error serializing orders: %v", err)
		return "", err
	}
	var orders []string
	var orderIds []string
	for _, order := range listOrder.Orders {
		for _, items := range order.OrderItems {
			product, err := rpc.ProductClient.GetProduct(ctx, &rpcproduct.GetProductReq{Id: items.Item.ProductId})
			if err != nil {
				klog.Errorf("Get product error: %v", err)
				return "", err
			}
			if lo.Contains(productsName, product.Product.Name) {
				jsonData, err = json.Marshal(order)
				if err != nil {
					klog.Errorf("Error serializing order: %v", err)
					return "", err
				}
				orders = append(orders, string(jsonData))
				orderIds = append(orderIds, order.OrderId)
				break
			}
		}
	}
	orderStr := strings.Join(orders, ",")
	// orderIdStr := strings.Join(orderIds, ";")
	return orderStr, nil
}

func searchOrdersTool() tool.InvokableTool {
	info := &schema.ToolInfo{
		Name: "获取指定商品的订单信息",
		Desc: "根据指定的商品名称，获取购买过指定商品的所有订单信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"names": {
				Desc: "用户指定的商品名称数组，例如：[手机, 电脑]",
				Type: schema.Array,
				ElemInfo: &schema.ParameterInfo{
					Type: schema.String, // 指定数组元素的类型为字符串
				},
				Required: true,
			},
		}),
	}

	return utils.NewTool(info, searchOrdersFunc)
}

func NewQueryOrderTools() []tool.BaseTool {
	// 初始化 tool
	tools := []tool.BaseTool{
		searchOrdersTool(),
	}
	return tools
}

func GetToolInfo(ctx context.Context, tools []tool.BaseTool) ([]*schema.ToolInfo, error) {
	toolInfos := make([]*schema.ToolInfo, 0, len(tools))
	for _, t := range tools {
		info, err := t.Info(ctx)
		if err != nil {
			klog.Errorf("get ToolInfo failed, err=%v", err)
			return nil, err
		}
		toolInfos = append(toolInfos, info)
	}
	return toolInfos, nil
}
