package service

import (
	"2501YTC/app/gateway/infra/rpc"
	"context"

	product "2501YTC/app/gateway/hertz_gen/gateway/product"
	rpcproduct "2501YTC/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type GetProductService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetProductService(ctx context.Context, requestContext *app.RequestContext) *GetProductService {
	return &GetProductService{RequestContext: requestContext, Context: ctx}
}

func (h *GetProductService) Run(req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	productResponse, err := rpc.ProductClient.GetProduct(h.Context, &rpcproduct.GetProductReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &product.GetProductResp{
		Product: &product.Product{
			Id:          productResponse.Product.Id,
			Picture:     productResponse.Product.Picture,
			Price:       productResponse.Product.Price,
			Description: productResponse.Product.Description,
			Name:        productResponse.Product.Name,
			Categories:  productResponse.Product.Categories,
		},
	}, nil
}
