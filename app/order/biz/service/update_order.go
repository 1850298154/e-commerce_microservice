package service

import (
	"context"
	"fmt"

	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	"2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/klog"
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
	if req.UserId == 0 || req.OrderId == "" {
		err = fmt.Errorf("user_id or order_id can not be empty")
		klog.Warnf("UpdateOrder failed, user_id or order_id can not be empty for Request %v", req)
		return
	}

	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 查询订单是否存在
		_, err = model.GetOrder(s.ctx, tx, req.UserId, req.OrderId)
		if err != nil {
			klog.Errorf("model.GetOrder.err:%v for UserId %v OrderId %v", err, req.UserId, req.OrderId)
			return err
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
			if err := model.UpdateOrder(s.ctx, tx, req.UserId, req.OrderId, updates); err != nil {
				return err
			}
		}

		// 更新订单项
		if len(req.NewOrderItems) > 0 {
			if err := model.UpdateOrderItems(s.ctx, tx, req.OrderId, req.NewOrderItems); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		klog.Errorf("UpdateOrder failed, UserId %v, OrderId %v err: %v", req.UserId, req.OrderId, err)
		return
	}

	return &order.UpdateOrderResp{Success: true}, nil
}
