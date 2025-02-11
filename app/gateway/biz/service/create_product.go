package service

import (
	"context"

	"2501YTC/app/gateway/infra/rpc"
	rpcproduct "2501YTC/rpc_gen/kitex_gen/product"

	product "2501YTC/app/gateway/hertz_gen/gateway/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type CreateProductService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateProductService(ctx context.Context, requestContext *app.RequestContext) *CreateProductService {
	return &CreateProductService{RequestContext: requestContext, Context: ctx}
}

func (h *CreateProductService) Run(req *product.CreateProductReq) (resp *product.CreateProductResp, err error) {
	_, err = rpc.ProductClient.CreateProduct(h.Context, &rpcproduct.CreateProductReq{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Picture:     req.Picture,
		Categories:  req.Categories,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
