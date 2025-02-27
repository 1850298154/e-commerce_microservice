package service

import (
	"context"

	rpcai "2501YTC/rpc_gen/kitex_gen/ai"

	"2501YTC/app/gateway/infra/rpc"
	"2501YTC/app/gateway/utils"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	ai "2501YTC/app/gateway/hertz_gen/gateway/ai"

	"github.com/cloudwego/hertz/pkg/app"
)

type QueryOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewQueryOrderService(Context context.Context, RequestContext *app.RequestContext) *QueryOrderService {
	return &QueryOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *QueryOrderService) Run(req *ai.QueryOrderReq) (resp *ai.QueryOrderResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	userId := utils.GetUserIdFromReqCtx(h.RequestContext)
	orderResp, err := rpc.AIClient.QueryOrder(h.Context, &rpcai.OrderQueryReq{
		UserId:  userId,
		Content: req.Content,
	})
	if err != nil {
		hlog.Errorf("auto order failed: %s", err.Error())
		return nil, err
	}

	orderItems := make([]*ai.OrderItem, 0)
	orderList := make([]*ai.OrderResult, 0)
	for _, order := range orderResp.Order {
		orderItems = make([]*ai.OrderItem, 0)
		for _, item := range order.OrderItems {
			orderItem := &ai.OrderItem{
				ProductId:   item.ProductId,
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				Cost:        item.Cost,
			}
			orderItems = append(orderItems, orderItem)
		}
		orderList = append(orderList, &ai.OrderResult{
			OrderId:      order.OrderId,
			UserId:       userId,
			UserCurrency: order.UserCurrency,
			Email:        order.Email,
			CreatedAt:    order.CreatedAt,
			OrderItems:   orderItems,
			OrderState:   order.OrderState,
		})
	}

	return &ai.QueryOrderResp{
		Order: orderList,
	}, nil
}
