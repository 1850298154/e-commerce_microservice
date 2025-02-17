package service

import (
	"context"

	"2501YTC/app/gateway/hertz_gen/gateway/user"
	"2501YTC/app/gateway/infra/rpc"

	rpcuser "2501YTC/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteUserService(ctx context.Context, requestContext *app.RequestContext) *DeleteUserService {
	return &DeleteUserService{RequestContext: requestContext, Context: ctx}
}

func (h *DeleteUserService) Run(req *user.DeleteUserReq) (resp *user.DeleteUserResp, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	res, err := rpc.UserClient.DeleteUser(h.Context, &rpcuser.DeleteUserReq{
		UserId: req.UserId,
	})
	if err != nil {
		return
	}
	return &user.DeleteUserResp{Success: res.Success}, nil
}
