package service

import (
	"context"
	"fmt"

	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	"2501YTC/rpc_gen/kitex_gen/cart"
	"2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/klog"
)

type ListOrderService struct {
	ctx context.Context
} // NewListOrderService new ListOrderService
func NewListOrderService(ctx context.Context) *ListOrderService {
	return &ListOrderService{ctx: ctx}
}

// Run 执行列出用户所有订单逻辑
func (s *ListOrderService) Run(req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	if req.UserId == 0 {
		err = fmt.Errorf("user id can not be empty")
		klog.Warn("ListOrder failed, UserId can not be empty")
		return
	}

	// 查询数据库获取订单信息
	orders, err := model.ListOrder(s.ctx, mysql.DB, req.UserId)
	if err != nil {
		klog.Errorf("model.ListOrder.err:%v for user id %v", err, req.UserId)
		return nil, err
	}

	// 将订单信息转换为rpc返回结构
	list := make([]*order.Order, 0, len(orders))
	for _, v := range orders {
		var items []*order.OrderItem
		for _, v := range v.OrderItems {
			items = append(items, &order.OrderItem{
				Cost: v.Cost,
				Item: &cart.CartItem{
					ProductId: v.ProductId,
					Quantity:  v.Quantity,
				},
			})
		}
		o := &order.Order{
			OrderId:      v.OrderId,
			UserId:       v.UserId,
			UserCurrency: v.UserCurrency,
			Email:        v.Consignee.Email,
			CreatedAt:    int32(v.CreatedAt.Unix()),
			Address: &order.Address{
				Country:       v.Consignee.Country,
				City:          v.Consignee.City,
				StreetAddress: v.Consignee.StreetAddress,
				ZipCode:       v.Consignee.ZipCode,
			},
			OrderItems: items,
		}
		list = append(list, o)
	}

	resp = &order.ListOrderResp{
		Orders: list,
	}
	return
}
