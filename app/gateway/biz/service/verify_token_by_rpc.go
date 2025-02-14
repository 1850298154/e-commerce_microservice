package service

import (
	"2501YTC/app/gateway/infra/rpc"
	rpcauth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"
	"errors"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"2501YTC/app/gateway/hertz_gen/gateway/auth"

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
	if req == nil {
		return nil, errors.New("请求对象为空")
	}
	token := req.Token
	if token == "" {
		return nil, errors.New("token为空")
	}
	hlog.CtxInfof(h.Context, "收到请求: %+v", req) // 使用结构化日志
	hlog.CtxDebugf(h.Context, "原始Token: %q", req.Token)
	if !strings.HasPrefix(token, "Bearer ") {
		return nil, errors.New("token缺少前缀Bearer")
	}
	token = token[len("Bearer "):]

	refreshToken := req.RefreshToken
	if refreshToken == "" {
		return nil, errors.New("refreshtoken为空")
	}
	if !strings.HasPrefix(refreshToken, "Bearer ") {
		return nil, errors.New("vrefreshtoken缺少前缀Bearer")
	}
	refreshToken = refreshToken[len("Bearer "):]
	rpcResponse, err := rpc.AuthClient.VerifyTokenByRPC(h.Context, &rpcauth.VerifyTokenReq{
		Token:        token,
		RefreshToken: refreshToken})
	if err != nil {
		return nil, err
	}
	return &auth.VerifyResp{
		Res: rpcResponse.Res,
	}, nil
}
