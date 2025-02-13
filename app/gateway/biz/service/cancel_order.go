package service

import (
	"context"

	"2501YTC/app/gateway/hertz_gen/gateway/order"
	"2501YTC/app/gateway/infra/rpc"

	rpcorder "2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/hertz/pkg/app"
)

type CancelOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCancelOrderService(ctx context.Context, requestContext *app.RequestContext) *CancelOrderService {
	return &CancelOrderService{RequestContext: requestContext, Context: ctx}
}

func (h *CancelOrderService) Run(req *order.CancelOrderReq) (resp *order.CancelOrderResp, err error) {
	rpcResponse, err := rpc.OrderClient.CancelOrder(h.Context, &rpcorder.CancelOrderReq{
		OrderId: req.OrderId,
		// TODO 从context获取UserId
		UserId:      req.UserId,
		TimedCancel: false,
		CancelTime:  req.CancelTime,
	})
	if err != nil {
		return nil, err
	}
	resp = &order.CancelOrderResp{
		Success: rpcResponse.Success,
	}

	return
}
