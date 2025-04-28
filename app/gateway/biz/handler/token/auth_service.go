package token

import (
	"context"

	"2501YTC/app/gateway/biz/service"
	"2501YTC/app/gateway/biz/utils"
	token "2501YTC/app/gateway/hertz_gen/gateway/token"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// DeliverTokenByRPC .
// @router /auth/token [POST]
func DeliverTokenByRPC(ctx context.Context, c *app.RequestContext) {
	var err error
	var req token.DeliverTokenReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &token.DeliveryResp{}
	resp, err = service.NewDeliverTokenByRPCService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// VerifyTokenByRPC .
// @router /auth/verify [POST]
func VerifyTokenByRPC(ctx context.Context, c *app.RequestContext) {
	var err error
	var req token.VerifyTokenReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &token.VerifyResp{}
	resp, err = service.NewVerifyTokenByRPCService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// RenewTokenByRPC .
// @router /auth/renew [POST]
func RenewTokenByRPC(ctx context.Context, c *app.RequestContext) {
	var err error
	var req token.RenewTokenReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &token.RenewTokenResp{}
	resp, err = service.NewRenewTokenByRPCService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
