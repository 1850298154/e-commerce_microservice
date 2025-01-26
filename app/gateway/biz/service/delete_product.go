package service

import (
	"context"

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
	// defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	// }()
	// todo edit your code
	return
}
