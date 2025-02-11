package service

import (
	"context"

	"2501YTC/app/gateway/infra/rpc"

	product "2501YTC/app/gateway/hertz_gen/gateway/product"
	rpcproduct "2501YTC/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type SearchProductsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSearchProductsService(ctx context.Context, requestContext *app.RequestContext) *SearchProductsService {
	return &SearchProductsService{RequestContext: requestContext, Context: ctx}
}

func (h *SearchProductsService) Run(req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	result, err := rpc.ProductClient.SearchProducts(h.Context, &rpcproduct.SearchProductsReq{
		Query:    req.Query,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &product.SearchProductsResp{
		Products: ConvertProductList(result.Results),
		Num:      result.Num,
	}, nil
}
