package service

import (
	"context"

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
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	return
}
