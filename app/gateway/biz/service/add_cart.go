package service

import (
	"context"

	cart "2501YTC/app/gateway/hertz_gen/gateway/cart"

	"github.com/cloudwego/hertz/pkg/app"
)

type AddCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddCartService(ctx context.Context, RequestContext *app.RequestContext) *AddCartService {
	return &AddCartService{RequestContext: RequestContext, Context: ctx}
}

func (h *AddCartService) Run(req *cart.AddCartReq) (resp *cart.Empty, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	return
}
