package service

import (
	"context"
	"errors"
	"strings"

	"2501YTC/app/gateway/hertz_gen/gateway/token"
	"2501YTC/app/gateway/infra/rpc"
	rpcauth "2501YTC/rpc_gen/kitex_gen/token"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/cloudwego/hertz/pkg/app"
)

type VerifyTokenByRPCService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewVerifyTokenByRPCService(ctx context.Context, requestContext *app.RequestContext) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{RequestContext: requestContext, Context: ctx}
}

func (h *VerifyTokenByRPCService) Run(req *token.VerifyTokenReq) (resp *token.VerifyResp, err error) {
	if req == nil {
		return nil, errors.New("请求对象为空")
	}
	t := req.Token
	if t == "" {
		return nil, errors.New("token为空")
	}
	hlog.CtxInfof(h.Context, "收到请求: %+v", req) // 使用结构化日志
	hlog.CtxDebugf(h.Context, "原始Token: %q", req.Token)
	if !strings.HasPrefix(t, "Bearer ") {
		return nil, errors.New("token缺少前缀Bearer")
	}
	t = t[len("Bearer "):]

	refreshToken := req.RefreshToken
	if refreshToken == "" {
		return nil, errors.New("refreshtoken为空")
	}
	if !strings.HasPrefix(refreshToken, "Bearer ") {
		return nil, errors.New("refreshtoken缺少前缀Bearer")
	}
	refreshToken = refreshToken[len("Bearer "):]
	rpcResponse, err := rpc.AuthClient.VerifyTokenByRPC(h.Context, &rpcauth.VerifyTokenReq{
		Token:        t,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}
	return &token.VerifyResp{
		Res: rpcResponse.Res,
	}, nil
}
