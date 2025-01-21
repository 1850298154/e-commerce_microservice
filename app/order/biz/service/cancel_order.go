package service

import (
	"context"
	"fmt"

	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	"2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/klog"
)

type CancelOrderService struct {
	ctx context.Context
} // NewCancelOrderService new CancelOrderService
func NewCancelOrderService(ctx context.Context) *CancelOrderService {
	return &CancelOrderService{ctx: ctx}
}

// Run 执行取消订单逻辑
func (s *CancelOrderService) Run(req *order.CancelOrderReq) (resp *order.CancelOrderResp, err error) {
	if req.UserId == 0 || req.OrderId == "" {
		err = fmt.Errorf("user_id or order_id can not be empty")
		klog.Warn("UpdateOrder failed, user_id or order_id can not be empty")
		return
	}

	// 查询订单是否存在
	_, err = model.GetOrder(s.ctx, mysql.DB, req.UserId, req.OrderId)
	if err != nil {
		klog.Errorf("model.GetOrder.err:%v for UserId %v OrderId %v", err, req.UserId, req.OrderId)
		return
	}
	if req.TimedCancel {
		err = model.CancelOrder(s.ctx, mysql.DB, req.UserId, req.OrderId, model.CancelTypeTimeout, req.CancelTime)
	} else {
		err = model.CancelOrder(s.ctx, mysql.DB, req.UserId, req.OrderId, model.CancelTypeUser, req.CancelTime)
	}

	if err != nil {
		klog.Errorf("model.CancelOrder.err:%v for UserId %v OrderId %v", err, req.UserId, req.OrderId)
		return
	}
	return &order.CancelOrderResp{Success: true}, nil
}
