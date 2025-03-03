package ai

import (
	ai "2501YTC/rpc_gen/kitex_gen/ai"
	"context"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func QueryOrder(ctx context.Context, req *ai.OrderQueryReq, callOptions ...callopt.Option) (resp *ai.OrderQueryResp, err error) {
	resp, err = defaultClient.QueryOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "QueryOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AutoOrder(ctx context.Context, req *ai.AutoOrderReq, callOptions ...callopt.Option) (resp *ai.AutoOrderResp, err error) {
	resp, err = defaultClient.AutoOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AutoOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
