package service

import (
	"context"

	product "2501YTC/app/gateway/hertz_gen/gateway/product"

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
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	return
}
