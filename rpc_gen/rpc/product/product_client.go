package product

import (
	product "2501YTC/rpc_gen/kitex_gen/product"
	"context"

	"2501YTC/rpc_gen/kitex_gen/product/productservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() productservice.Client
	Service() string
	CreateProduct(ctx context.Context, Req *product.CreateProductReq, callOptions ...callopt.Option) (r *product.CreateProductResp, err error)
	UpdateProduct(ctx context.Context, Req *product.UpdateProductReq, callOptions ...callopt.Option) (r *product.UpdateProductResp, err error)
	DeleteProduct(ctx context.Context, Req *product.DeleteProductReq, callOptions ...callopt.Option) (r *product.DeleteProductResp, err error)
	ListProducts(ctx context.Context, Req *product.ListProductsReq, callOptions ...callopt.Option) (r *product.ListProductsResp, err error)
	GetProduct(ctx context.Context, Req *product.GetProductReq, callOptions ...callopt.Option) (r *product.GetProductResp, err error)
	SearchProducts(ctx context.Context, Req *product.SearchProductsReq, callOptions ...callopt.Option) (r *product.SearchProductsResp, err error)
	SearchProductsByName(ctx context.Context, Req *product.SearchProductsReq, callOptions ...callopt.Option) (r *product.SearchProductsResp, err error)
	UploadImage(ctx context.Context, Req *product.UploadImageReq, callOptions ...callopt.Option) (r *product.UploadImageResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := productservice.NewClient(dstService, opts...)
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
	kitexClient productservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() productservice.Client {
	return c.kitexClient
}

func (c *clientImpl) CreateProduct(ctx context.Context, Req *product.CreateProductReq, callOptions ...callopt.Option) (r *product.CreateProductResp, err error) {
	return c.kitexClient.CreateProduct(ctx, Req, callOptions...)
}

func (c *clientImpl) UpdateProduct(ctx context.Context, Req *product.UpdateProductReq, callOptions ...callopt.Option) (r *product.UpdateProductResp, err error) {
	return c.kitexClient.UpdateProduct(ctx, Req, callOptions...)
}

func (c *clientImpl) DeleteProduct(ctx context.Context, Req *product.DeleteProductReq, callOptions ...callopt.Option) (r *product.DeleteProductResp, err error) {
	return c.kitexClient.DeleteProduct(ctx, Req, callOptions...)
}

func (c *clientImpl) ListProducts(ctx context.Context, Req *product.ListProductsReq, callOptions ...callopt.Option) (r *product.ListProductsResp, err error) {
	return c.kitexClient.ListProducts(ctx, Req, callOptions...)
}

func (c *clientImpl) GetProduct(ctx context.Context, Req *product.GetProductReq, callOptions ...callopt.Option) (r *product.GetProductResp, err error) {
	return c.kitexClient.GetProduct(ctx, Req, callOptions...)
}

func (c *clientImpl) SearchProducts(ctx context.Context, Req *product.SearchProductsReq, callOptions ...callopt.Option) (r *product.SearchProductsResp, err error) {
	return c.kitexClient.SearchProducts(ctx, Req, callOptions...)
}

func (c *clientImpl) SearchProductsByName(ctx context.Context, Req *product.SearchProductsReq, callOptions ...callopt.Option) (r *product.SearchProductsResp, err error) {
	return c.kitexClient.SearchProductsByName(ctx, Req, callOptions...)
}

func (c *clientImpl) UploadImage(ctx context.Context, Req *product.UploadImageReq, callOptions ...callopt.Option) (r *product.UploadImageResp, err error) {
	return c.kitexClient.UploadImage(ctx, Req, callOptions...)
}
