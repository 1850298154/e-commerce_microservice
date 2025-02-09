package service

import (
	"context"

	"2501YTC/app/gateway/infra/rpc"
	rpcproduct "2501YTC/rpc_gen/kitex_gen/product"

	product "2501YTC/app/gateway/hertz_gen/gateway/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteProductService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteProductService(ctx context.Context, requestContext *app.RequestContext) *DeleteProductService {
	return &DeleteProductService{RequestContext: requestContext, Context: ctx}
}

func (h *DeleteProductService) Run(req *product.DeleteProductReq) (resp *product.DeleteProductResp, err error) {
	_, err = rpc.ProductClient.DeleteProduct(h.Context, &rpcproduct.DeleteProductReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
