package service

import (
	"context"

	auth "2501YTC/app/gateway/hertz_gen/gateway/auth"

	"github.com/cloudwego/hertz/pkg/app"
)

type RenewTokenByRPCService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRenewTokenByRPCService(ctx context.Context, requestContext *app.RequestContext) *RenewTokenByRPCService {
	return &RenewTokenByRPCService{RequestContext: requestContext, Context: ctx}
}

func (h *RenewTokenByRPCService) Run(req *auth.RenewTokenReq) (resp *auth.RenewTokenResp, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	return
}
