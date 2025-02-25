package service

import (
	"context"

	"2501YTC/app/gateway/hertz_gen/gateway/user"
	"2501YTC/app/gateway/infra/rpc"
	"2501YTC/app/gateway/utils"

	rpcuser "2501YTC/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
)

type GetUserInfoService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserInfoService(ctx context.Context, requestContext *app.RequestContext) *GetUserInfoService {
	return &GetUserInfoService{RequestContext: requestContext, Context: ctx}
}

func (h *GetUserInfoService) Run(req *user.GetUserInfoReq) (resp *user.GetUserInfoResp, err error) {
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	userId := utils.GetUserIdFromReqCtx(h.RequestContext)
	res, err := rpc.UserClient.GetUserInfo(h.Context, &rpcuser.GetUserInfoReq{
		UserId: userId,
	})
	if err != nil {
		return
	}
	return &user.GetUserInfoResp{
		UserId:    res.UserId,
		Email:     res.Email,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}
