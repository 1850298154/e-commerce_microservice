package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	auth "2501YTC/app/gateway/hertz_gen/gateway/auth"
	"2501YTC/app/gateway/infra/rpc"
	rpcauth "2501YTC/rpc_gen/kitex_gen/auth"

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

	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(h.Context, time.Second*5)
	defer cancel()

	rpcResponse, err := rpc.AuthClient.DeliverTokenByRPC(ctx, &rpcauth.DeliverTokenReq{UserId: req.UserId, Role: 2})
	if err != nil {
		return nil, fmt.Errorf("RPC调用失败:%v", err)
	}
	if rpcResponse == nil {
		return nil, errors.New("RPC返回空响应")
	}
	h.RequestContext.Response.Header.Set("Authorization", "Bearer "+rpcResponse.Token)
	h.RequestContext.Response.Header.Set("X-Refresh-Token", "Bearer "+rpcResponse.RefreshToken)
	//	h.RequestContext.Response.Header.Set("Set-Cookie", "Authorization=Bearer "+rpcResponse.Token+"; Path=/; HttpOnly; Secure; SameSite=Lax")
	return &auth.DeliveryResp{
		Token:        rpcResponse.Token,
		RefreshToken: rpcResponse.RefreshToken,
	}, nil
}
