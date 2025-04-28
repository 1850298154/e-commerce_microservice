package token

import (
	token "2501YTC/rpc_gen/kitex_gen/token"
	"context"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func DeliverTokenByRPC(ctx context.Context, req *token.DeliverTokenReq, callOptions ...callopt.Option) (resp *token.DeliveryResp, err error) {
	resp, err = defaultClient.DeliverTokenByRPC(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeliverTokenByRPC call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func VerifyTokenByRPC(ctx context.Context, req *token.VerifyTokenReq, callOptions ...callopt.Option) (resp *token.VerifyResp, err error) {
	resp, err = defaultClient.VerifyTokenByRPC(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "VerifyTokenByRPC call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func RenewTokenByRPC(ctx context.Context, req *token.RenewTokenReq, callOptions ...callopt.Option) (resp *token.RenewTokenResp, err error) {
	resp, err = defaultClient.RenewTokenByRPC(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "RenewTokenByRPC call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteTokenByRPC(ctx context.Context, req *token.DeleteTokenReq, callOptions ...callopt.Option) (resp *token.DeleteTokenResp, err error) {
	resp, err = defaultClient.DeleteTokenByRPC(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteTokenByRPC call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
