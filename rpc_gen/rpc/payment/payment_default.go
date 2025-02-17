package payment

import (
	payment "2501YTC/rpc_gen/kitex_gen/payment"
	"context"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func Charge(ctx context.Context, req *payment.ChargeReq, callOptions ...callopt.Option) (resp *payment.ChargeResp, err error) {
	resp, err = defaultClient.Charge(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Charge call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func CancelPayment(ctx context.Context, req *payment.CancelPaymentReq, callOptions ...callopt.Option) (resp *payment.CancelPaymentResp, err error) {
	resp, err = defaultClient.CancelPayment(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CancelPayment call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func TimedCancelPayment(ctx context.Context, req *payment.TimedCancelPaymentReq, callOptions ...callopt.Option) (resp *payment.TimedCancelPaymentResp, err error) {
	resp, err = defaultClient.TimedCancelPayment(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "TimedCancelPayment call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
