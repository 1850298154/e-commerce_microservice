package service

import (
	"context"

	"2501YTC/app/gateway/infra/rpc"
	"2501YTC/app/gateway/utils"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	ai "2501YTC/app/gateway/hertz_gen/gateway/ai"
	rpcai "2501YTC/rpc_gen/kitex_gen/ai"

	"github.com/cloudwego/hertz/pkg/app"
)

type AutoOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAutoOrderService(ctx context.Context, requestContext *app.RequestContext) *AutoOrderService {
	return &AutoOrderService{RequestContext: requestContext, Context: ctx}
}

func (h *AutoOrderService) Run(req *ai.PlaceOrderReq) (resp *ai.PlaceOrderResp, err error) {
	// defer func() {
	//  hlog.CtxInfof(h.Context, "req = %+v", req)
	//  hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	userId := utils.GetUserIdFromReqCtx(h.RequestContext)
	orderResp, err := rpc.AIClient.AutoOrder(h.Context, &rpcai.AutoOrderReq{
		UserId:  userId,
		Content: req.Content,
	})
	if err != nil {
		hlog.Errorf("auto order failed: %s", err.Error())
		return nil, err
	}

	orderItems := make([]*ai.OrderItem, 0)
	for _, item := range orderResp.Order.OrderItems {
		orderItem := &ai.OrderItem{
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Cost:        item.Cost,
		}
		orderItems = append(orderItems, orderItem)
	}

	order := &ai.OrderResult{
		OrderId:      orderResp.Order.OrderId,
		UserId:       userId,
		UserCurrency: orderResp.Order.UserCurrency,
		Email:        orderResp.Order.Email,
		CreatedAt:    orderResp.Order.CreatedAt,
		OrderItems:   orderItems,
		OrderState:   orderResp.Order.OrderState,
	}
	return &ai.PlaceOrderResp{
		Order: order,
	}, nil
}
