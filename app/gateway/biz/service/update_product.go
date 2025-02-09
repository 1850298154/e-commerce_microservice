package service

import (
	"context"

	"2501YTC/app/gateway/infra/rpc"
	rpcproduct "2501YTC/rpc_gen/kitex_gen/product"

	product "2501YTC/app/gateway/hertz_gen/gateway/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateProductService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateProductService(ctx context.Context, requestContext *app.RequestContext) *UpdateProductService {
	return &UpdateProductService{RequestContext: requestContext, Context: ctx}
}

func (h *UpdateProductService) Run(req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	_, err = rpc.ProductClient.UpdateProduct(h.Context, &rpcproduct.UpdateProductReq{
		Id:          req.Id,
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
