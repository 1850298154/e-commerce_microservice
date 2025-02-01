package service

import (
	"context"

	"2501YTC/app/gateway/hertz_gen/gateway/user"
	"2501YTC/app/gateway/infra/rpc"
	rpcuser "2501YTC/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateUserService(Context context.Context, RequestContext *app.RequestContext) *UpdateUserService {
	return &UpdateUserService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateUserService) Run(req *user.UpdateUserReq) (resp *user.UpdateUserResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	res, err := rpc.UserClient.UpdateUser(h.Context, &rpcuser.UpdateUserReq{
		UserId:   req.UserId,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return
	}
	return &user.UpdateUserResp{Success: res.Success}, nil
}
