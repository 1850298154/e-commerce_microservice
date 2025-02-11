package service

import (
	"context"

	"2501YTC/app/gateway/infra/rpc"

	product "2501YTC/app/gateway/hertz_gen/gateway/product"
	rpcproduct "2501YTC/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type ListProductsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListProductsService(ctx context.Context, requestContext *app.RequestContext) *ListProductsService {
	return &ListProductsService{RequestContext: requestContext, Context: ctx}
}

func (h *ListProductsService) Run(req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	list, err := rpc.ProductClient.ListProducts(h.Context, &rpcproduct.ListProductsReq{
		Page:         int32(req.Page),
		PageSize:     int64(req.Size),
		CategoryName: req.CategoryName,
	})
	if err != nil {
		return nil, err
	}
	return &product.ListProductsResp{
		Products: ConvertProductList(list.Products),
		Num:      uint32(list.Num),
	}, nil
}

func ConvertProductList(list []*rpcproduct.Product) []*product.Product {
	products := make([]*product.Product, 0, len(list))
	for _, p := range list {
		products = append(products, &product.Product{
			Id:          p.Id,
			Picture:     p.Picture,
			Price:       p.Price,
			Description: p.Description,
			Name:        p.Name,
			Categories:  p.Categories,
		})
	}
	return products
}
