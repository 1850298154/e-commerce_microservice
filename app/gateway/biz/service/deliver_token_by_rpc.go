package service

import (
	"2501YTC/app/gateway/infra/rpc"
	rpcauth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"

	auth "2501YTC/app/gateway/hertz_gen/gateway/auth"

	"github.com/cloudwego/hertz/pkg/app"
)

type DeliverTokenByRPCService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeliverTokenByRPCService(ctx context.Context, requestContext *app.RequestContext) *DeliverTokenByRPCService {
	return &DeliverTokenByRPCService{RequestContext: requestContext, Context: ctx}
}

func (h *DeliverTokenByRPCService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	rpcResponse, err := rpc.AuthClient.DeliverTokenByRPC(h.Context, &rpcauth.DeliverTokenReq{UserId: req.UserId})
	return &auth.DeliveryResp{
			Token:        rpcResponse.Token,
			RefreshToken: rpcResponse.RefreshToken},
		nil
}
