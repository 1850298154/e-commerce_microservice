package service

import (
	"context"

	product "2501YTC/app/gateway/hertz_gen/gateway/product"

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
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	return
}
