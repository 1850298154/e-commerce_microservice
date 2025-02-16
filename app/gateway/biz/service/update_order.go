package service

import (
	"context"

	"2501YTC/app/gateway/hertz_gen/gateway/order"
	"2501YTC/app/gateway/infra/rpc"
	"2501YTC/app/gateway/utils"
	rpccart "2501YTC/rpc_gen/kitex_gen/cart"
	rpcorder "2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateOrderService(ctx context.Context, requestContext *app.RequestContext) *UpdateOrderService {
	return &UpdateOrderService{RequestContext: requestContext, Context: ctx}
}

func (h *UpdateOrderService) Run(req *order.UpdateOrderReq) (resp *order.UpdateOrderResp, err error) {
	newOrderItems := make([]*rpcorder.OrderItem, 0, len(req.NewOrderItems))
	for _, item := range req.NewOrderItems {
		newOrderItems = append(newOrderItems, &rpcorder.OrderItem{
			Item: &rpccart.CartItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			},
			Cost: item.Cost,
		})
	}
	rpcResponse, err := rpc.OrderClient.UpdateOrder(h.Context, &rpcorder.UpdateOrderReq{
		OrderId: req.OrderId,
		// TODO 从context获取UserId
		// UserId: req.UserId,
		UserId: utils.GetUserIdFromReqCtx(h.RequestContext),
		NewAddress: &rpcorder.Address{
			StreetAddress: req.NewAddress.StreetAddress,
			City:          req.NewAddress.City,
			State:         req.NewAddress.State,
			Country:       req.NewAddress.Country,
			ZipCode:       req.NewAddress.ZipCode,
		},
		NewOrderItems: newOrderItems,
		NewEmail:      req.NewEmail,
	})
	if err != nil {
		return nil, err
	}
	return &order.UpdateOrderResp{
		Success: rpcResponse.Success,
	}, nil
}
