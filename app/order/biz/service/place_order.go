package service

import (
	"context"
	"database/sql"
	"fmt"

	"2501YTC/app/order/biz/dal/mq"
	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	Error "2501YTC/app/order/error"
	"2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
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
	// TODO tracing place order
	_, span := otel.Tracer("order server").Start(s.ctx, "PlaceOrderService.Run")
	defer span.End()

	// 参数校验
	if len(req.OrderItems) == 0 {
		err = fmt.Errorf("order items empty")
		klog.Warnf("PlaceOrder failed, OrderItems empty, UserId: %d", req.UserId)
		return nil, Error.NewError(Error.ErrInvalidOrderItems, "order items empty", nil)
	}
	if req.UserId == 0 {
		err = fmt.Errorf("user id can not be empty")
		klog.CtxWarnf(s.ctx, "PlaceOrder failed, UserId can not be empty")
		return nil, Error.NewError(Error.ErrInvalidUserId, "user id can not be empty", nil)
	}

	// 开启事务
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		orderQuery := model.NewOrderQuery(s.ctx, tx)
		// 生成订单ID
		orderId, err := uuid.NewUUID()
		if err != nil {
			klog.CtxErrorf(s.ctx, "PlaceOrder failed, generate order id failed, UserId: %d, err: %v", req.UserId, err)
			return Error.NewError(Error.ErrGenerateOrderIdFailed, "generate order id failed", err)
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

		// 设置订单商品
		if len(req.OrderItems) != 0 {
			o.OrderItems = make([]model.OrderItem, len(req.OrderItems))
			for idx, v := range req.OrderItems {
				o.OrderItems[idx] = model.OrderItem{
					OrderIdRefer: o.OrderId,
					ProductId:    v.Item.ProductId,
					Quantity:     v.Item.Quantity,
					Cost:         v.Cost,
				}
			}
		}
		// 创建订单
		if err := orderQuery.CreateOrder(o); err != nil {
			klog.CtxErrorf(s.ctx, "PlaceOrder failed, create order failed, UserId: %d, err: %v", req.UserId, err)
			return Error.NewError(Error.ErrCreateOrderFailed, "create order failed", err)
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
				klog.CtxErrorf(s.ctx, "PlaceOrder failed, send delay message failed, UserId: %d, err: %v", req.UserId, err)
			}
		}()
		return nil
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "PlaceOrder failed, UserId: %d, err: %v", req.UserId, err)
	}
	klog.CtxInfof(s.ctx, "PlaceOrder success for UserId %v as orderId %v", req.UserId, resp.Order.OrderId)
	return resp, err
}
