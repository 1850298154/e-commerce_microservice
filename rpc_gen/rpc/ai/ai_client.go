package ai

import (
	ai "2501YTC/rpc_gen/kitex_gen/ai"
	"context"

	"2501YTC/rpc_gen/kitex_gen/ai/authservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() authservice.Client
	Service() string
	SearchforOrder(ctx context.Context, Req *ai.SearchforOrderReq, callOptions ...callopt.Option) (r *ai.SearchforOrderResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := authservice.NewClient(dstService, opts...)
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
	kitexClient authservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() authservice.Client {
	return c.kitexClient
}

func (c *clientImpl) SearchforOrder(ctx context.Context, Req *ai.SearchforOrderReq, callOptions ...callopt.Option) (r *ai.SearchforOrderResp, err error) {
	return c.kitexClient.SearchforOrder(ctx, Req, callOptions...)
}
