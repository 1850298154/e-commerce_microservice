package order

import (
	"context"

	"2501YTC/app/gateway/biz/service"
	"2501YTC/app/gateway/biz/utils"
	order "2501YTC/app/gateway/hertz_gen/gateway/order"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// PlaceOrder .
// @router /orders [POST]
func PlaceOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req order.PlaceOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &order.PlaceOrderResp{}
	resp, err = service.NewPlaceOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ListOrder .
// @router /orders [GET]
func ListOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req order.ListOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &order.ListOrderResp{}
	resp, err = service.NewListOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// MarkOrderPaid .
// @router /orders/{order_id}/paid [PUT]
func MarkOrderPaid(ctx context.Context, c *app.RequestContext) {
	var err error
	var req order.MarkOrderPaidReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &order.MarkOrderPaidResp{}
	resp, err = service.NewMarkOrderPaidService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UpdateOrder .
// @router /orders/{order_id} [PUT]
func UpdateOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req order.UpdateOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &order.UpdateOrderResp{}
	resp, err = service.NewUpdateOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// CancelOrder .
// @router /orders/{order_id} [DELETE]
func CancelOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req order.CancelOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &order.CancelOrderResp{}
	resp, err = service.NewCancelOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
