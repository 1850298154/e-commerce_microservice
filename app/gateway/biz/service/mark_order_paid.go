package service

import (
	"context"

	"2501YTC/app/gateway/hertz_gen/gateway/order"
	"2501YTC/app/gateway/infra/rpc"
	"2501YTC/app/gateway/utils"
	rpcorder "2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/hertz/pkg/app"
)

type MarkOrderPaidService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewMarkOrderPaidService(ctx context.Context, requestContext *app.RequestContext) *MarkOrderPaidService {
	return &MarkOrderPaidService{RequestContext: requestContext, Context: ctx}
}

func (h *MarkOrderPaidService) Run(req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	rpcResponse, err := rpc.OrderClient.MarkOrderPaid(h.Context, &rpcorder.MarkOrderPaidReq{
		OrderId: req.OrderId,
		// TODO 从context获取UserId
		// UserId: req.UserId,
		UserId: utils.GetUserIdFromReqCtx(h.RequestContext),
	})
	if err != nil {
		return nil, err
	}
	return &order.MarkOrderPaidResp{
		Success: rpcResponse.Success,
	}, nil
}
