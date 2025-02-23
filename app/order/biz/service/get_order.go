package service

import (
	"context"
	"fmt"

	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/biz/model"
	Error "2501YTC/app/order/error"
	"2501YTC/rpc_gen/kitex_gen/cart"
	"2501YTC/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/klog"
	"go.opentelemetry.io/otel"
)

type GetOrderService struct {
	ctx context.Context
} // NewGetOrderService new GetOrderService
func NewGetOrderService(ctx context.Context) *GetOrderService {
	return &GetOrderService{ctx: ctx}
}

// Run 执行获取订单逻辑
func (s *GetOrderService) Run(req *order.GetOrderReq) (resp *order.GetOrderResp, err error) {
	// tracing get order
	_, span := otel.Tracer("order server").Start(s.ctx, "GetOrderService.Run")
	defer span.End()

	if req.UserId == 0 || req.OrderId == "" {
		err = fmt.Errorf("user_id or order_id can not be empty")
		klog.CtxWarnf(s.ctx, "GetOrder failed, user_id or order_id can not be empty")
		return nil, Error.NewError(Error.ErrInvalidUserId, "user_id or order_id can not be empty", nil)
	}

	orderQuery := model.NewOrderQuery(s.ctx, mysql.DB)
	// 查询订单是否存在
	curOrder, err := orderQuery.GetOrder(req.OrderId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "model.GetOrder.err:%v for UserId %v OrderId %v", err, req.UserId, req.OrderId)
		return nil, Error.NewError(Error.ErrGetOrderByUserIdAndOrderIdFailed, fmt.Sprintf("GetOrder failed for UserId %v OrderId %v", req.UserId, req.OrderId), err)
	}

	if curOrder.UserId != req.UserId {
		klog.CtxWarnf(s.ctx, "GetOrder failed, UserId %v does not match OrderId %v", req.UserId, req.OrderId)
		return nil, Error.NewError(Error.ErrGetOrderByUserIdAndOrderIdFailed, fmt.Sprintf("GetOrder failed, UserId %v does not match OrderId %v", req.UserId, req.OrderId), nil)
	}

	res := &order.GetOrderResp{
		Order: &order.OrderResp{
			Order: &order.Order{
				OrderId:      curOrder.OrderId,
				UserId:       curOrder.UserId,
				UserCurrency: curOrder.UserCurrency,
				Email:        curOrder.Consignee.Email,
				CreatedAt:    int32(curOrder.CreatedAt.Unix()),
				Address: &order.Address{
					StreetAddress: curOrder.Consignee.StreetAddress,
					City:          curOrder.Consignee.City,
					State:         curOrder.Consignee.State,
					Country:       curOrder.Consignee.Country,
					ZipCode:       curOrder.Consignee.ZipCode,
				},
			},
			OrderState: string(curOrder.OrderState),
		},
		Success: true,
	}

	if len(curOrder.OrderItems) > 0 {
		res.Order.Order.OrderItems = make([]*order.OrderItem, len(curOrder.OrderItems))
		for idx, item := range curOrder.OrderItems {
			res.Order.Order.OrderItems[idx] = &order.OrderItem{
				Item: &cart.CartItem{
					ProductId: item.ProductId,
					Quantity:  item.Quantity,
				},
				Cost: item.Cost,
			}
		}
	}
	klog.CtxInfof(s.ctx, "GetOrder success for UserId %v OrderId %v, \nOrder %v", req.UserId, req.OrderId, res.Order)
	return res, nil
}
