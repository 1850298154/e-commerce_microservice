package user

import (
	user "2501YTC/rpc_gen/kitex_gen/user"
	"context"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func Register(ctx context.Context, req *user.RegisterReq, callOptions ...callopt.Option) (resp *user.RegisterResp, err error) {
	resp, err = defaultClient.Register(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Register call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func Login(ctx context.Context, req *user.LoginReq, callOptions ...callopt.Option) (resp *user.LoginResp, err error) {
	resp, err = defaultClient.Login(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Login call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteUser(ctx context.Context, req *user.DeleteUserReq, callOptions ...callopt.Option) (resp *user.DeleteUserResp, err error) {
	resp, err = defaultClient.DeleteUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateUser(ctx context.Context, req *user.UpdateUserReq, callOptions ...callopt.Option) (resp *user.UpdateUserResp, err error) {
	resp, err = defaultClient.UpdateUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetUserInfo(ctx context.Context, req *user.GetUserInfoReq, callOptions ...callopt.Option) (resp *user.GetUserInfoResp, err error) {
	resp, err = defaultClient.GetUserInfo(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetUserInfo call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateUserRole(ctx context.Context, req *user.UpdateUserRoleReq, callOptions ...callopt.Option) (resp *user.UpdateUserRoleResp, err error) {
	resp, err = defaultClient.UpdateUserRole(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateUserRole call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
