package service

import (
	"context"

	auth "2501YTC/app/gateway/hertz_gen/gateway/auth"

	"github.com/cloudwego/hertz/pkg/app"
)

type VerifyTokenByRPCService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewVerifyTokenByRPCService(ctx context.Context, requestContext *app.RequestContext) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{RequestContext: requestContext, Context: ctx}
}

func (h *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	return
}
