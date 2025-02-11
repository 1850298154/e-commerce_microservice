package service

import (
	"2501YTC/app/gateway/hertz_gen/gateway/order"
	"2501YTC/app/gateway/infra/rpc"
	"context"

	rpccart "2501YTC/rpc_gen/kitex_gen/cart"
	rpcorder "2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/hertz/pkg/app"
)

type PlaceOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewPlaceOrderService(ctx context.Context, requestContext *app.RequestContext) *PlaceOrderService {
	return &PlaceOrderService{RequestContext: requestContext, Context: ctx}
}

func (h *PlaceOrderService) Run(req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	newOrderItems := make([]*rpcorder.OrderItem, 0, len(req.OrderItems))
	for _, item := range req.OrderItems {
		newOrderItems = append(newOrderItems, &rpcorder.OrderItem{
			Item: &rpccart.CartItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			},
			Cost: item.Cost,
		})
	}
	rpcResponse, err := rpc.OrderClient.PlaceOrder(h.Context, &rpcorder.PlaceOrderReq{
		// TODO 从context获取UserId
		UserId:       req.UserId,
		UserCurrency: req.UserCurrency,
		Address: &rpcorder.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       req.Address.ZipCode,
		},
		Email:      req.Email,
		OrderItems: newOrderItems,
	})
	if err != nil {
		return nil, err
	}
	return &order.PlaceOrderResp{
		OrderId: rpcResponse.Order.OrderId,
	}, nil
}
