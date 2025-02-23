package service

import (
	"context"

	"2501YTC/app/gateway/infra/rpc"
	"2501YTC/app/gateway/utils"

	cart "2501YTC/app/gateway/hertz_gen/gateway/cart"

	rpccart "2501YTC/rpc_gen/kitex_gen/cart"

	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteCartService(ctx context.Context, requestContext *app.RequestContext) *DeleteCartService {
	return &DeleteCartService{RequestContext: requestContext, Context: ctx}
}

func (h *DeleteCartService) Run(req *cart.Empty) (resp *cart.Empty, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	_, err = rpc.CartClient.EmptyCart(h.Context, &rpccart.EmptyCartReq{
		UserId: utils.GetUserIdFromReqCtx(h.RequestContext),
		// UserId: 1,
	})
	if err != nil {
		return nil, err
	}
	return
}
