package ai

import (
	ai "2501YTC/rpc_gen/kitex_gen/ai"
	"context"

	"2501YTC/rpc_gen/kitex_gen/ai/orderservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() orderservice.Client
	Service() string
	QueryOrder(ctx context.Context, Req *ai.OrderQueryReq, callOptions ...callopt.Option) (r *ai.OrderQueryResp, err error)
	AutoOrder(ctx context.Context, Req *ai.AutoOrderReq, callOptions ...callopt.Option) (r *ai.AutoOrderResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := orderservice.NewClient(dstService, opts...)
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
	kitexClient orderservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() orderservice.Client {
	return c.kitexClient
}

func (c *clientImpl) QueryOrder(ctx context.Context, Req *ai.OrderQueryReq, callOptions ...callopt.Option) (r *ai.OrderQueryResp, err error) {
	return c.kitexClient.QueryOrder(ctx, Req, callOptions...)
}

func (c *clientImpl) AutoOrder(ctx context.Context, Req *ai.AutoOrderReq, callOptions ...callopt.Option) (r *ai.AutoOrderResp, err error) {
	return c.kitexClient.AutoOrder(ctx, Req, callOptions...)
}
