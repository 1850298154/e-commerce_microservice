package auth

import (
	"2501YTC/app/gateway/biz/service"
	"2501YTC/app/gateway/biz/utils"
	auth "2501YTC/app/gateway/hertz_gen/gateway/auth"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// DeliverTokenByRPC .
// @router /auth/token [POST]
func DeliverTokenByRPC(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.DeliverTokenReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	// fmt.Println("Received UserId:", req.UserId)

	resp := &auth.DeliveryResp{}
	resp, err = service.NewDeliverTokenByRPCService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// VerifyTokenByRPC .
// @router /auth/renew [POST]
func VerifyTokenByRPC(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.VerifyTokenReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &auth.VerifyResp{}
	resp, err = service.NewVerifyTokenByRPCService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// RenewTokenByRPC .
// @router /auth/verify [POST]
func RenewTokenByRPC(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.RenewTokenReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &auth.RenewTokenResp{}
	resp, err = service.NewRenewTokenByRPCService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
