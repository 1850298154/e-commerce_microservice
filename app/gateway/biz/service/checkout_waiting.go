package service

import (
	"context"

	common "2501YTC/app/gateway/hertz_gen/gateway/common"

	"github.com/cloudwego/hertz/pkg/app"
)

type CheckoutWaitingService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCheckoutWaitingService(ctx context.Context, requestContext *app.RequestContext) *CheckoutWaitingService {
	return &CheckoutWaitingService{RequestContext: requestContext, Context: ctx}
}

func (h *CheckoutWaitingService) Run(req *common.Empty) (resp *common.Empty, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	return
}
