package service

import (
	"context"

	"2501YTC/app/gateway/hertz_gen/gateway/order"
	"2501YTC/app/gateway/infra/rpc"
	rpcorder "2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/hertz/pkg/app"
)

type ListOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListOrderService(ctx context.Context, requestContext *app.RequestContext) *ListOrderService {
	return &ListOrderService{RequestContext: requestContext, Context: ctx}
}

func (h *ListOrderService) Run(req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	rpcResponse, err := rpc.OrderClient.ListOrder(h.Context, &rpcorder.ListOrderReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	orders := make([]*order.Order, 0, len(rpcResponse.Orders))
	for _, order_ := range rpcResponse.Orders {
		orders = append(orders, &order.Order{
			OrderId:      order_.OrderId,
			UserId:       order_.UserId,
			UserCurrency: order_.UserCurrency,
			Address: &order.Address{
				StreetAddress: order_.Address.StreetAddress,
				City:          order_.Address.City,
				State:         order_.Address.State,
				Country:       order_.Address.Country,
				ZipCode:       order_.Address.ZipCode,
			},
			Email:     order_.Email,
			CreatedAt: order_.CreatedAt,
			OrderItems: func() []*order.OrderItem {
				items := make([]*order.OrderItem, 0, len(order_.OrderItems))
				for _, item := range order_.OrderItems {
					items = append(items, &order.OrderItem{
						ProductId: item.Item.ProductId,
						Quantity:  item.Item.Quantity,
						Cost:      item.Cost,
					})
				}
				return items
			}(),
		})
	}
	return &order.ListOrderResp{
		Orders: orders,
	}, nil
}
