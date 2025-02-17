package payment

import (
	payment "2501YTC/rpc_gen/kitex_gen/payment"
	"context"

	"2501YTC/rpc_gen/kitex_gen/payment/paymentservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() paymentservice.Client
	Service() string
	Charge(ctx context.Context, Req *payment.ChargeReq, callOptions ...callopt.Option) (r *payment.ChargeResp, err error)
	CancelPayment(ctx context.Context, Req *payment.CancelPaymentReq, callOptions ...callopt.Option) (r *payment.CancelPaymentResp, err error)
	TimedCancelPayment(ctx context.Context, Req *payment.TimedCancelPaymentReq, callOptions ...callopt.Option) (r *payment.TimedCancelPaymentResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := paymentservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient paymentservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() paymentservice.Client {
	return c.kitexClient
}

func (c *clientImpl) Charge(ctx context.Context, Req *payment.ChargeReq, callOptions ...callopt.Option) (r *payment.ChargeResp, err error) {
	return c.kitexClient.Charge(ctx, Req, callOptions...)
}

func (c *clientImpl) CancelPayment(ctx context.Context, Req *payment.CancelPaymentReq, callOptions ...callopt.Option) (r *payment.CancelPaymentResp, err error) {
	return c.kitexClient.CancelPayment(ctx, Req, callOptions...)
}

func (c *clientImpl) TimedCancelPayment(ctx context.Context, Req *payment.TimedCancelPaymentReq, callOptions ...callopt.Option) (r *payment.TimedCancelPaymentResp, err error) {
	return c.kitexClient.TimedCancelPayment(ctx, Req, callOptions...)
}
