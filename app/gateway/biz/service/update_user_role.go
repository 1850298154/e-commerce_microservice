package service

import (
	"context"

	"2501YTC/app/gateway/hertz_gen/gateway/user"
	"2501YTC/app/gateway/infra/rpc"
	rpcuser "2501YTC/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateUserRoleService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateUserRoleService(Context context.Context, RequestContext *app.RequestContext) *UpdateUserRoleService {
	return &UpdateUserRoleService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateUserRoleService) Run(req *user.UpdateUserRoleReq) (resp *user.UpdateUserRoleResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	res, err := rpc.UserClient.UpdateUserRole(h.Context, &rpcuser.UpdateUserRoleReq{
		UserId: req.UserId,
		Role:   req.Role,
	})
	if err != nil {
		return
	}
	return &user.UpdateUserRoleResp{Success: res.Success}, nil
}
