package service

import (
	"context"
	"fmt"

	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	Error "2501YTC/app/order/error"
	"2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/klog"
	"go.opentelemetry.io/otel"
)

type CancelOrderService struct {
	ctx context.Context
} // NewCancelOrderService new CancelOrderService
func NewCancelOrderService(ctx context.Context) *CancelOrderService {
	return &CancelOrderService{ctx: ctx}
}

// Run 执行取消订单逻辑
func (s *CancelOrderService) Run(req *order.CancelOrderReq) (resp *order.CancelOrderResp, err error) {
	// TODO tracing cancel order
	_, span := otel.Tracer("order server").Start(s.ctx, "CancelOrderService.Run")
	defer span.End()

	if req.UserId == 0 || req.OrderId == "" {
		err = fmt.Errorf("user_id or order_id can not be empty")
		// klog.Warnf( "UpdateOrder failed, user_id or order_id can not be empty")

		klog.CtxWarnf(s.ctx, "UpdateOrder failed, user_id or order_id can not be empty")
		return nil, Error.NewError(Error.ErrInvalidUserId, "user_id or order_id can not be empty", nil)
	}

	orderQuery := model.NewOrderQuery(s.ctx, mysql.DB)
	// 查询订单是否存在
	curOrder, err := orderQuery.GetOrder(req.OrderId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "model.GetOrder.err:%v for UserId %v OrderId %v", err, req.UserId, req.OrderId)
		return nil, Error.NewError(Error.ErrGetOrderByUserIdAndOrderIdFailed, fmt.Sprintf("GetOrder failed for UserId %v OrderId %v", req.UserId, req.OrderId), err)
	}

	// 订单已取消，不允许再次取消
	if curOrder.OrderState == model.OrderStateCanceled {
		klog.CtxWarnf(s.ctx, "CancelOrder failed, OrderId %v has been canceled", req.OrderId)
		return nil, Error.NewError(Error.ErrCancelOrderFailed, fmt.Sprintf("CancelOrder failed, OrderId %v has been canceled", req.OrderId), nil)
	}
	if curOrder.UserId != req.UserId {
		klog.CtxWarnf(s.ctx, "CancelOrder failed, UserId %v does not match OrderId %v", req.UserId, req.OrderId)
		return nil, Error.NewError(Error.ErrCancelOrderFailed, fmt.Sprintf("CancelOrder failed, UserId %v does not match OrderId %v", req.UserId, req.OrderId), nil)
	}

	if req.TimedCancel {
		err = orderQuery.CancelOrder(req.OrderId, model.CancelTypeTimeout, req.CancelTime)
	} else {
		err = orderQuery.CancelOrder(req.OrderId, model.CancelTypeUser, req.CancelTime)
	}

	if err != nil {
		klog.CtxErrorf(s.ctx, "model.CancelOrder.err:%v for UserId %v OrderId %v", err, req.UserId, req.OrderId)
		return nil, Error.NewError(Error.ErrCancelOrderFailed, fmt.Sprintf("CancelOrder failed for UserId %v OrderId %v", req.UserId, req.OrderId), err)
	}
	klog.CtxInfof(s.ctx, "CancelOrder success for UserId %v OrderId %v", req.UserId, req.OrderId)
	return &order.CancelOrderResp{Success: true}, nil
}
