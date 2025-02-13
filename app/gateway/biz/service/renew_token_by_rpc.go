package service

import (
	"context"

	auth "2501YTC/app/gateway/hertz_gen/gateway/auth"
	"2501YTC/app/gateway/infra/rpc"
	rpcauth "2501YTC/rpc_gen/kitex_gen/auth"

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
	rpcResponse, err := rpc.AuthClient.RenewTokenByRPC(h.Context, &rpcauth.RenewTokenReq{RefreshToken: req.RefreshToken})
	return &auth.RenewTokenResp{
		Token:        rpcResponse.Token,
		RefreshToken: rpcResponse.RefreshToken,
		ExpiresIn:    rpcResponse.ExpiresIn,
	}, nil
}
