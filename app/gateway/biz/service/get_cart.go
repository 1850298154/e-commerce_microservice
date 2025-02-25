package service

import (
	"context"

	"2501YTC/app/gateway/biz/utils"
	"2501YTC/app/gateway/infra/rpc"
	uutils "2501YTC/app/gateway/utils"

	cart "2501YTC/app/gateway/hertz_gen/gateway/cart"

	rpccart "2501YTC/rpc_gen/kitex_gen/cart"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type GetCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetCartService(ctx context.Context, requestContext *app.RequestContext) *GetCartService {
	return &GetCartService{RequestContext: requestContext, Context: ctx}
}

func (h *GetCartService) Run(req *cart.Empty) (resp *cart.Empty, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code

	getcartresp, err := rpc.CartClient.GetCart(h.Context, &rpccart.GetCartReq{
		UserId: uutils.GetUserIdFromReqCtx(h.RequestContext),
		// UserId: 1,
	})
	if err != nil {
		return nil, err
	}
	utils.SendSuccessResponse(h.Context, h.RequestContext, consts.StatusOK, getcartresp)
	return nil, nil
}
