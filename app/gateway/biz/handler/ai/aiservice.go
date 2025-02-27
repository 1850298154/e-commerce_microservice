package ai

import (
	"context"

	"2501YTC/app/gateway/biz/service"
	"2501YTC/app/gateway/biz/utils"
	ai "2501YTC/app/gateway/hertz_gen/gateway/ai"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// QueryOrder .
// @router /ai/query [POST]
func QueryOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req ai.QueryOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	c.Set("user_id", uint32(1))

	resp, err := service.NewQueryOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// AutoOrder .
// @router /ai/place [POST]
func AutoOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req ai.PlaceOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	c.Set("user_id", uint32(1))

	resp, err := service.NewAutoOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
