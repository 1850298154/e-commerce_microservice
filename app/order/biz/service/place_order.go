package service

import (
	"context"
	"fmt"

	"2501YTC/app/order/biz/dal/mq"
	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	"2501YTC/rpc_gen/kitex_gen/order"

	"database/sql"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
		err = fmt.Errorf("order items empty")
		klog.Warnf("PlaceOrder failed, OrderItems empty, UserId: %d", req.UserId)
		return
	}
	if req.UserId == 0 {
		err = fmt.Errorf("user id can not be empty")
		klog.Warn("PlaceOrder failed, UserId can not be empty")
		return
	}

	// 开启事务
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 生成订单ID
		orderId, err := uuid.NewUUID()
		if err != nil {
			klog.Errorf("PlaceOrder failed, generate order id failed, UserId: %d, err: %v", req.UserId, err)
		}

		o := &model.Order{
			OrderId:      orderId.String(),
			OrderState:   model.OrderStatePlaced,
			UserId:       req.UserId,
			UserCurrency: req.UserCurrency,
			Consignee: model.Consignee{
				Email: req.Email,
			},
			CancelTime: sql.NullTime{Valid: false}, // 显式设置为 NULL
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
			klog.Errorf("PlaceOrder failed, create order failed, UserId: %d, err: %v", req.UserId, err)
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
			klog.Errorf("PlaceOrder failed, create order item failed, UserId: %d, err: %v", req.UserId, err)
			return err
		}

		resp = &order.PlaceOrderResp{
			Order: &order.OrderResult{
				OrderId: orderId.String(),
			},
		}

		// 发送定时取消订单消息到消息队列
		go func() {
			err := mq.ProducerInstance.SendDelayMessage(orderId.String())
			if err != nil {
				klog.Errorf("PlaceOrder failed, send delay message failed, UserId: %d, err: %v", req.UserId, err)
			}
		}()
		return nil
	})
	if err != nil {
		klog.Errorf("PlaceOrder failed, UserId: %d, err: %v", req.UserId, err)
	}

	return
}
