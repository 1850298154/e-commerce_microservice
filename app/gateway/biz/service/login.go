package service

import (
	"context"
	"fmt"

	"2501YTC/app/gateway/hertz_gen/gateway/user"
	"2501YTC/app/gateway/infra/rpc"

	rpcauth "2501YTC/rpc_gen/kitex_gen/token"
	rpcuser "2501YTC/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
)

type LoginService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLoginService(ctx context.Context, requestContext *app.RequestContext) *LoginService {
	return &LoginService{RequestContext: requestContext, Context: ctx}
}

func (h *LoginService) Run(req *user.LoginReq) (resp *user.LoginResp, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	res, err := rpc.UserClient.Login(h.Context, &rpcuser.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return
	}
	fmt.Println(res.UserId)
	tokenReq, err := rpc.AuthClient.DeliverTokenByRPC(h.Context, &rpcauth.DeliverTokenReq{
		UserId: res.UserId,
		Role:   res.Role,
	})
	if err != nil {
		return
	}
	h.RequestContext.Response.Header.Set("Authorization", "Bearer "+tokenReq.Token)
	h.RequestContext.Response.Header.Set("X-Refresh-Token", tokenReq.RefreshToken)

	return &user.LoginResp{
		UserId:       res.UserId,
		Token:        tokenReq.Token,
		RefreshToken: tokenReq.RefreshToken,
	}, nil
}
