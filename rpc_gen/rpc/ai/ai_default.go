package ai

import (
	ai "2501YTC/rpc_gen/kitex_gen/ai"
	"context"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func SearchforOrder(ctx context.Context, req *ai.SearchforOrderReq, callOptions ...callopt.Option) (resp *ai.SearchforOrderResp, err error) {
	resp, err = defaultClient.SearchforOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SearchforOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
