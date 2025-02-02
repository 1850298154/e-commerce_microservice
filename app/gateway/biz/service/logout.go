package service

import (
	"2501YTC/app/gateway/infra/rpc"
	"2501YTC/rpc_gen/kitex_gen/auth"
	"context"
	"errors"
	"strings"

	"2501YTC/app/gateway/hertz_gen/gateway/user"

	"github.com/cloudwego/hertz/pkg/app"
)

type LogoutService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLogoutService(ctx context.Context, requestContext *app.RequestContext) *LogoutService {
	return &LogoutService{RequestContext: requestContext, Context: ctx}
}

func (h *LogoutService) Run(req *user.LogoutReq) (resp *user.LogoutResp, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()

	token := req.Token
	if token == "" {
		return nil, errors.New("token为空")
	}
	if !strings.HasPrefix(token, "Bearer ") {
		return nil, errors.New("token含前缀bearer")
	}
	token = token[len("Bearer "):]
	_, err = rpc.AuthClient.DeleteTokenByRPC(h.Context, &auth.DeleteTokenReq{
		Token: token,
	})
	if err != nil {
		return nil, err
	}
	return &user.LogoutResp{
		Success: true,
	}, nil
}
