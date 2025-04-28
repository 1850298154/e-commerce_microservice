package token

import (
	token "2501YTC/rpc_gen/kitex_gen/token"
	"context"

	"2501YTC/rpc_gen/kitex_gen/token/authservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() authservice.Client
	Service() string
	DeliverTokenByRPC(ctx context.Context, Req *token.DeliverTokenReq, callOptions ...callopt.Option) (r *token.DeliveryResp, err error)
	VerifyTokenByRPC(ctx context.Context, Req *token.VerifyTokenReq, callOptions ...callopt.Option) (r *token.VerifyResp, err error)
	RenewTokenByRPC(ctx context.Context, Req *token.RenewTokenReq, callOptions ...callopt.Option) (r *token.RenewTokenResp, err error)
	DeleteTokenByRPC(ctx context.Context, Req *token.DeleteTokenReq, callOptions ...callopt.Option) (r *token.DeleteTokenResp, err error)
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

func (c *clientImpl) DeliverTokenByRPC(ctx context.Context, Req *token.DeliverTokenReq, callOptions ...callopt.Option) (r *token.DeliveryResp, err error) {
	return c.kitexClient.DeliverTokenByRPC(ctx, Req, callOptions...)
}

func (c *clientImpl) VerifyTokenByRPC(ctx context.Context, Req *token.VerifyTokenReq, callOptions ...callopt.Option) (r *token.VerifyResp, err error) {
	return c.kitexClient.VerifyTokenByRPC(ctx, Req, callOptions...)
}

func (c *clientImpl) RenewTokenByRPC(ctx context.Context, Req *token.RenewTokenReq, callOptions ...callopt.Option) (r *token.RenewTokenResp, err error) {
	return c.kitexClient.RenewTokenByRPC(ctx, Req, callOptions...)
}

func (c *clientImpl) DeleteTokenByRPC(ctx context.Context, Req *token.DeleteTokenReq, callOptions ...callopt.Option) (r *token.DeleteTokenResp, err error) {
	return c.kitexClient.DeleteTokenByRPC(ctx, Req, callOptions...)
}
