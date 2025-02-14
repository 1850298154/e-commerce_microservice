package service

import (
	"2501YTC/app/gateway/infra/rpc"
	rpcauth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"
	"errors"
	"fmt"
	"strings"

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
	refreshToken := req.RefreshToken
	if refreshToken == "" {
		return nil, errors.New("refreshtoken为空")
	}
	if !strings.HasPrefix(refreshToken, "Bearer ") {
		return nil, errors.New("refreshtoken缺少前缀Bearer")
	}
	refreshToken = refreshToken[len("Bearer "):]

	rpcResponse, err := rpc.AuthClient.RenewTokenByRPC(h.Context, &rpcauth.RenewTokenReq{RefreshToken: refreshToken})
	if err != nil {
		return nil, fmt.Errorf("renewtokenRPC调用失败: %w", err)
	}
	// if rpcResponse == nil {
	//	return nil, errors.New("renewtoken响应为空")
	// }
	if rpcResponse.Token == "" || rpcResponse.RefreshToken == "" {
		return nil, errors.New("令牌生成异常")
	}
	return &auth.RenewTokenResp{
		Token:        rpcResponse.Token,
		RefreshToken: rpcResponse.RefreshToken,
		ExpiresIn:    rpcResponse.ExpiresIn,
	}, nil
}
