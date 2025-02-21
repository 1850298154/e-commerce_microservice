package service

import (
	"context"

	"2501YTC/app/gateway/hertz_gen/gateway/order"
	"2501YTC/app/gateway/infra/rpc"
	"2501YTC/app/gateway/utils"
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
		// UserId: req.UserId,
		UserId: utils.GetUserIdFromReqCtx(h.RequestContext),
	})
	if err != nil {
		return nil, err
	}
	orders := make([]*order.OrderResp, 0, len(rpcResponse.Orders))
	for _, order_ := range rpcResponse.Orders {
		orders = append(orders, &order.OrderResp{
			Order: &order.Order{
				OrderId:      order_.Order.OrderId,
				UserId:       order_.Order.UserId,
				UserCurrency: order_.Order.UserCurrency,
				Address: &order.Address{
					StreetAddress: order_.Order.Address.StreetAddress,
					City:          order_.Order.Address.City,
					State:         order_.Order.Address.State,
					Country:       order_.Order.Address.Country,
					ZipCode:       order_.Order.Address.ZipCode,
				},
				Email:     order_.Order.Email,
				CreatedAt: order_.Order.CreatedAt,
				OrderItems: func() []*order.OrderItem {
					items := make([]*order.OrderItem, 0, len(order_.Order.OrderItems))
					for _, item := range order_.Order.OrderItems {
						items = append(items, &order.OrderItem{
							ProductId: item.Item.ProductId,
							Quantity:  item.Item.Quantity,
							Cost:      item.Cost,
						})
					}
					return items
				}(),
			},
			OrderState: order_.OrderState,
		})
	}
	return &order.ListOrderResp{
		Orders: orders,
	}, nil
}
