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

type MarkOrderPaidService struct {
	ctx context.Context
} // NewMarkOrderPaidService new MarkOrderPaidService
func NewMarkOrderPaidService(ctx context.Context) *MarkOrderPaidService {
	return &MarkOrderPaidService{ctx: ctx}
}

// Run 执行标记订单支付逻辑
func (s *MarkOrderPaidService) Run(req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	// TODO tracing mark order paid
	_, span := otel.Tracer("order server").Start(s.ctx, "MarkOrderPaidService.Run")
	defer span.End()

	if req.UserId == 0 || req.OrderId == "" {
		err = fmt.Errorf("user_id or order_id can not be empty")
		klog.CtxWarnf(s.ctx, "MarkOrderPaid failed, user_id or order_id can not be empty")
		return nil, Error.NewError(Error.ErrInvalidUserId, "user_id or order_id can not be empty", nil)
	}

	orderQuery := model.NewOrderQuery(s.ctx, mysql.DB)
	// 查询要更新的订单
	curOrder, err := orderQuery.GetOrder(req.OrderId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "model.GetOrder.err:%v for UserId %v OrderId %v", err, req.UserId, req.OrderId)
		return nil, Error.NewError(Error.ErrGetOrderByUserIdAndOrderIdFailed, fmt.Sprintf("GetOrder failed for UserId %v OrderId %v", req.UserId, req.OrderId), err)
	}

	// 订单已取消或已支付，不允许再次支付
	if curOrder.OrderState != model.OrderStatePlaced {
		klog.CtxWarnf(s.ctx, "MarkOrderPaid failed, OrderId %v state is not placed", req.OrderId)
		return nil, Error.NewError(Error.ErrUpdateOrderFailed, fmt.Sprintf("MarkOrderPaid failed, OrderId %v state is not placed", req.OrderId), nil)
	}
	if curOrder.UserId != req.UserId {
		klog.CtxWarnf(s.ctx, "MarkOrderPaid failed, UserId %v does not match OrderId %v", req.UserId, req.OrderId)
		return nil, Error.NewError(Error.ErrUpdateOrderFailed, fmt.Sprintf("MarkOrderPaid failed, UserId %v does not match OrderId %v", req.UserId, req.OrderId), nil)
	}

	// 更新订单状态为已支付
	err = orderQuery.UpdateOrderState(req.OrderId, model.OrderStatePaid)
	if err != nil {
		klog.CtxErrorf(s.ctx, "model.UpdateOrderState.err:%v for UserId %v OrderId %v", err, req.UserId, req.OrderId)
		return nil, Error.NewError(Error.ErrUpdateOrderFailed, fmt.Sprintf("UpdateOrderState failed for UserId %v OrderId %v", req.UserId, req.OrderId), err)
	}

	resp = &order.MarkOrderPaidResp{
		Success: true,
	}
	klog.CtxInfof(s.ctx, "MarkOrderPaid success for UserId %v OrderId %v", req.UserId, req.OrderId)
	return
}
