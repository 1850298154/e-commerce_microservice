package service

import (
	"context"
	"fmt"

	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	"2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type PlaceOrderService struct {
	ctx context.Context
}

// NewPlaceOrderService 新建PlaceOrderService
func NewPlaceOrderService(ctx context.Context) *PlaceOrderService {
	return &PlaceOrderService{ctx: ctx}
}

// Run 执行创建订单逻辑
func (s *PlaceOrderService) Run(req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	if len(req.OrderItems) == 0 {
		err = fmt.Errorf("OrderItems empty")
		klog.Warn("PlaceOrder failed, OrderItems empty, UserId: %d", req.UserId)
		return
	}

	// 开启事务
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 生成订单ID
		orderId, err := uuid.NewUUID()
		if err != nil {
			klog.Error("PlaceOrder failed, generate order id failed, UserId: %d, err: %v", req.UserId)
		}

		o := &model.Order{
			OrderId:      orderId.String(),
			OrderState:   model.OrderStatePlaced,
			UserId:       req.UserId,
			UserCurrency: req.UserCurrency,
			Consignee: model.Consignee{
				Email: req.Email,
			},
		}

		// 设置收货地址
		if req.Address != nil {
			a := req.Address
			o.Consignee.Country = a.Country
			o.Consignee.State = a.State
			o.Consignee.City = a.City
			o.Consignee.StreetAddress = a.StreetAddress
		}

		// 创建订单
		if err := tx.Create(o).Error; err != nil {
			klog.Error("PlaceOrder failed, create order failed, UserId: %d, err: %v", req.UserId, err)
			return err
		}

		// 创建订单商品项
		var itemList []*model.OrderItem
		for _, v := range req.OrderItems {
			itemList = append(itemList, &model.OrderItem{
				OrderIdRefer: o.OrderId,
				ProductId:    v.Item.ProductId,
				Quantity:     v.Item.Quantity,
				Cost:         v.Cost,
			})
		}
		if err := tx.Create(&itemList).Error; err != nil {
			klog.Error("PlaceOrder failed, create order item failed, UserId: %d, err: %v", req.UserId, err)
			return err
		}

		resp = &order.PlaceOrderResp{
			Order: &order.OrderResult{
				OrderId: orderId.String(),
			},
		}
		return nil
	})

	if err != nil {
		klog.Error("PlaceOrder failed, UserId: %d, err: %v", req.UserId, err)
	}

	return
}
