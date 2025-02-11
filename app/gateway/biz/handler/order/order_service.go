package order

import (
	"2501YTC/app/gateway/biz/service"
	"2501YTC/app/gateway/biz/utils"
	"context"

	hertzorder "2501YTC/app/gateway/hertz_gen/gateway/order"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	ErrBindAndValidateFailed = 400
	ErrRPCCallFailed         = 500
	RPCCallSuccess           = 200
)

// PlaceOrder .
// @router /orders [POST]
func PlaceOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertzorder.PlaceOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Warnf("bind and validate failed, %v", err)
		utils.SendErrResponse(ctx, c, ErrBindAndValidateFailed, err)
		return
	}

	resp := &hertzorder.PlaceOrderResp{}
	resp, err = service.NewPlaceOrderService(ctx, c).Run(&req)
	if err != nil {
		hlog.Warnf("place order failed, %v", err)
		utils.SendErrResponse(ctx, c, ErrRPCCallFailed, err)
		return
	}
	hlog.Infof("place order success, order_id: %d", resp.OrderId)
	utils.SendSuccessResponse(ctx, c, RPCCallSuccess, resp)
}

// ListOrder .
// @router /orders [GET]
func ListOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertzorder.ListOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, ErrBindAndValidateFailed, err)
		hlog.Warnf("bind and validate failed, %v", err)
		return
	}

	resp := &hertzorder.ListOrderResp{}
	resp, err = service.NewListOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, ErrRPCCallFailed, err)
		hlog.Warnf("list order failed, %v", err)
		return
	}
	hlog.Infof("list order success, order count: %d", len(resp.Orders))
	utils.SendSuccessResponse(ctx, c, RPCCallSuccess, resp)
}

// MarkOrderPaid .
// @router /orders/paid [POST]
func MarkOrderPaid(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertzorder.MarkOrderPaidReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, ErrBindAndValidateFailed, err)
		hlog.Warnf("bind and validate failed, %v", err)
		return
	}

	resp := &hertzorder.MarkOrderPaidResp{}
	resp, err = service.NewMarkOrderPaidService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, ErrRPCCallFailed, err)
		hlog.Warnf("mark order paid failed, %v", err)
		return
	}
	hlog.Infof("mark order paid success, order_id: %d", req.OrderId)
	utils.SendSuccessResponse(ctx, c, RPCCallSuccess, resp)
}

// UpdateOrder .
// @router /orders [PUT]
func UpdateOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertzorder.UpdateOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, ErrBindAndValidateFailed, err)
		return
	}

	resp := &hertzorder.UpdateOrderResp{}
	resp, err = service.NewUpdateOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, ErrRPCCallFailed, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, RPCCallSuccess, resp)
}

// CancelOrder .
// @router /orders [DELETE]
func CancelOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertzorder.CancelOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, ErrBindAndValidateFailed, err)
		return
	}

	resp := &hertzorder.CancelOrderResp{}
	resp, err = service.NewCancelOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, ErrRPCCallFailed, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, RPCCallSuccess, resp)
}
