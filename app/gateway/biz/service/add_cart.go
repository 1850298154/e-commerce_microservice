package service

import (
	"context"

	cart "2501YTC/app/gateway/hertz_gen/gateway/cart"
	"2501YTC/app/gateway/infra/rpc"
	"2501YTC/app/gateway/utils"

	rpccart "2501YTC/rpc_gen/kitex_gen/cart"

	"github.com/cloudwego/hertz/pkg/app"
)

type AddCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddCartService(ctx context.Context, requestContext *app.RequestContext) *AddCartService {
	return &AddCartService{RequestContext: requestContext, Context: ctx}
}

func (h *AddCartService) Run(req *cart.AddCartReq) (resp *cart.Empty, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	_, err = rpc.CartClient.AddItem(h.Context, &rpccart.AddItemReq{
		UserId: utils.GetUserIdFromReqCtx(h.RequestContext),
		// UserId: 1,
		Item: &rpccart.CartItem{
			ProductId: req.ProductId,
			Quantity:  req.ProductNum,
		},
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
