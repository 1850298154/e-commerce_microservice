package service

import (
	"context"

	auth "2501YTC/app/gateway/hertz_gen/gateway/auth"

	"github.com/cloudwego/hertz/pkg/app"
)

type DeliverTokenByRPCService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeliverTokenByRPCService(Context context.Context, RequestContext *app.RequestContext) *DeliverTokenByRPCService {
	return &DeliverTokenByRPCService{RequestContext: RequestContext, Context: Context}
}

func (h *DeliverTokenByRPCService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
