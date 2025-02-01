package service

import (
	"2501YTC/app/gateway/hertz_gen/gateway/auth"
	"2501YTC/app/gateway/infra/rpc"
	"context"
	rpcauth "2501YTC/rpc_gen/kitex_gen/auth"
	rpcuser "2501YTC/rpc_gen/kitex_gen/user"

	"2501YTC/app/gateway/hertz_gen/gateway/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type RegisterService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRegisterService(Context context.Context, RequestContext *app.RequestContext) *RegisterService {
	return &RegisterService{RequestContext: RequestContext, Context: Context}
}

func (h *RegisterService) Run(req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	res, err := rpc.UserClient.Register(h.Context, &rpcuser.RegisterReq{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		return nil, err
	}

	token, err := rpc.AuthClient.DeliverTokenByRPC(h.Context, &rpcauth.DeliverTokenReq{
		UserId: res.UserId,
		Role:   res.Role,
	})
	return
}
