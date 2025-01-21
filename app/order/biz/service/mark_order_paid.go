package service

import (
	"context"
	"fmt"

	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	"2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/klog"
)

type MarkOrderPaidService struct {
	ctx context.Context
} // NewMarkOrderPaidService new MarkOrderPaidService
func NewMarkOrderPaidService(ctx context.Context) *MarkOrderPaidService {
	return &MarkOrderPaidService{ctx: ctx}
}

// Run 执行标记订单支付逻辑
func (s *MarkOrderPaidService) Run(req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	if req.UserId == 0 || req.OrderId == "" {
		err = fmt.Errorf("user_id or order_id can not be empty")
		klog.Warn("MarkOrderPaid failed, user_id or order_id can not be empty")
		return
	}

	// 查询要更新的订单
	_, err = model.GetOrder(s.ctx, mysql.DB, req.UserId, req.OrderId)
	if err != nil {
		klog.Errorf("model.GetOrder.err:%v for UserId %v OrderId %v", err, req.UserId, req.OrderId)
		return
	}

	// 更新订单状态为已支付
	err = model.UpdateOrderState(s.ctx, mysql.DB, req.UserId, req.OrderId, model.OrderStatePaid)
	if err != nil {
		klog.Errorf("model.UpdateOrderState.err:%v for UserId %v OrderId %v", err, req.UserId, req.OrderId)
		return
	}

	resp = &order.MarkOrderPaidResp{}
	return
}
