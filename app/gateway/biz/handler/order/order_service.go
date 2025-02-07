package order

import (
	"context"

	"2501YTC/app/gateway/biz/service"
	"2501YTC/app/gateway/biz/utils"
	hertzorder "2501YTC/app/gateway/hertz_gen/gateway/order"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	ErrBindAndValidateFailed = 4000
	PlaceOrderFailed         = 4081
	ListOrderFailed          = 4082
	MarkOrderPaidFailed      = 4083
	UpdateOrderFailed        = 4084
	CancelOrderFailed        = 4085

	PlaceOrderSuccess    = 2081
	ListOrderSuccess     = 2082
	MarkOrderPaidSuccess = 2083
	UpdateOrderSuccess   = 2084
	CancelOrderSuccess   = 2085
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
		utils.SendErrResponse(ctx, c, PlaceOrderFailed, err)
		return
	}
	hlog.Infof("place order success, order_id: %d", resp.OrderId)
	utils.SendSuccessResponse(ctx, c, PlaceOrderSuccess, resp)
}

// ListOrder .
// @router /orders [GET]
func ListOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertzorder.ListOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, ErrBindAndValidateFailed, err)
		return
	}

	resp := &hertzorder.ListOrderResp{}
	resp, err = service.NewListOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, ListOrderFailed, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, ListOrderSuccess, resp)
}

// MarkOrderPaid .
// @router /orders/{order_id}/paid [PUT]
func MarkOrderPaid(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertzorder.MarkOrderPaidReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, ErrBindAndValidateFailed, err)
		return
	}

	resp := &hertzorder.MarkOrderPaidResp{}
	resp, err = service.NewMarkOrderPaidService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, MarkOrderPaidFailed, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, MarkOrderPaidSuccess, resp)
}

// UpdateOrder .
// @router /orders/{order_id} [PUT]
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
		utils.SendErrResponse(ctx, c, UpdateOrderFailed, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, UpdateOrderSuccess, resp)
}

// CancelOrder .
// @router /orders/{order_id} [DELETE]
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
		utils.SendErrResponse(ctx, c, CancelOrderFailed, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, CancelOrderSuccess, resp)
}
