package service

import (
	"context"

	cart "2501YTC/app/gateway/hertz_gen/gateway/cart"

	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteCartService(ctx context.Context, RequestContext *app.RequestContext) *DeleteCartService {
	return &DeleteCartService{RequestContext: RequestContext, Context: ctx}
}

func (h *DeleteCartService) Run(req *cart.Empty) (resp *cart.Empty, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	return
}
