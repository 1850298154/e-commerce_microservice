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
	"gorm.io/gorm"
)

type UpdateOrderService struct {
	ctx context.Context
} // NewUpdateOrderService new UpdateOrderService
func NewUpdateOrderService(ctx context.Context) *UpdateOrderService {
	return &UpdateOrderService{ctx: ctx}
}

// Run 执行更新订单信息逻辑
func (s *UpdateOrderService) Run(req *order.UpdateOrderReq) (resp *order.UpdateOrderResp, err error) {
	// tracing update order
	_, span := otel.Tracer("order server").Start(s.ctx, "UpdateOrderService.Run")
	defer span.End()

	if req.UserId == 0 || req.OrderId == "" {
		// err = fmt.Errorf("user_id or order_id can not be empty")
		err = Error.NewError(Error.ErrInvalidUserId, "user_id or order_id can not be empty", nil)
		klog.CtxWarnf(s.ctx, "UpdateOrder failed, user_id or order_id can not be empty for Request %v", req)
		return
	}

	if req.NewEmail == "" && req.NewAddress == nil && len(req.NewOrderItems) == 0 {
		// err = fmt.Errorf("no field to update")
		err = Error.NewError(Error.ErrUpdateOrderFailed, "no field to update", nil)
		klog.CtxWarnf(s.ctx, "UpdateOrder failed, no field to update for Request %v", req)
		return
	}

	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		orderQuery := model.NewOrderQuery(s.ctx, tx)
		// 查询订单是否存在
		curOrder, err := orderQuery.GetOrder(req.OrderId)
		if err != nil {
			klog.CtxErrorf(s.ctx, "model.GetOrder.err:%v for UserId %v OrderId %v", err, req.UserId, req.OrderId)
			return Error.NewError(Error.ErrGetOrderByUserIdAndOrderIdFailed, fmt.Sprintf("GetOrder failed for UserId %v OrderId %v", req.UserId, req.OrderId), err)
		}

		// 订单已取消，不允许更新
		if curOrder.OrderState == model.OrderStateCanceled {
			klog.CtxWarnf(s.ctx, "UpdateOrder failed, OrderId %v has been canceled", req.OrderId)
			return Error.NewError(Error.ErrUpdateOrderFailed, fmt.Sprintf("UpdateOrder failed, OrderId %v has been canceled", req.OrderId), nil)
		}
		if curOrder.UserId != req.UserId {
			klog.CtxWarnf(s.ctx, "UpdateOrder failed, UserId %v does not match OrderId %v", req.UserId, req.OrderId)
			return Error.NewError(Error.ErrUpdateOrderFailed, fmt.Sprintf("UpdateOrder failed, UserId %v does not match OrderId %v", req.UserId, req.OrderId), nil)
		}

		// 处理更新字段
		updates := map[string]any{}
		if req.NewEmail != "" {
			updates["email"] = req.NewEmail
		}
		if req.NewAddress != nil {
			updates["city"] = req.NewAddress.City
			updates["state"] = req.NewAddress.State
			updates["country"] = req.NewAddress.Country
			updates["zip_code"] = req.NewAddress.ZipCode
			updates["street_address"] = req.NewAddress.StreetAddress
		}

		// 更新订单基本信息
		if len(updates) > 0 {
			if err := orderQuery.UpdateOrder(req.OrderId, updates); err != nil {
				klog.CtxErrorf(s.ctx, "UpdateOrder failed for UserId %v OrderId %v", req.UserId, req.OrderId)
				return Error.NewError(Error.ErrUpdateOrderFailed, fmt.Sprintf("UpdateOrder failed for UserId %v OrderId %v", req.UserId, req.OrderId), err)
			}
		}

		// 更新订单项
		if len(req.NewOrderItems) > 0 {
			if err := orderQuery.UpdateOrderItems(req.OrderId, req.NewOrderItems); err != nil {
				klog.CtxErrorf(s.ctx, "UpdateOrderItems failed for OrderId %v", req.OrderId)
				return Error.NewError(Error.ErrUpdateOrderItemsFailed, fmt.Sprintf("UpdateOrderItems failed for OrderId %v", req.OrderId), err)
			}
		}

		return nil
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "UpdateOrder failed, UserId %v, OrderId %v err: %v", req.UserId, req.OrderId, err)
		return nil, err
	}
	klog.CtxInfof(s.ctx, "UpdateOrder success for UserId %v OrderId %v", req.UserId, req.OrderId)
	return &order.UpdateOrderResp{Success: true}, nil
}
