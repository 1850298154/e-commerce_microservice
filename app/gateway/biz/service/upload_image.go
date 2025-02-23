package service

import (
	"context"

	"2501YTC/app/gateway/infra/rpc"

	rpcproduct "2501YTC/rpc_gen/kitex_gen/product"

	product "2501YTC/app/gateway/hertz_gen/gateway/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type UploadImageService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUploadImageService(ctx context.Context, requestContext *app.RequestContext) *UploadImageService {
	return &UploadImageService{RequestContext: requestContext, Context: ctx}
}

func (h *UploadImageService) Run(req *product.UploadImageReq) (resp *product.UploadImageResp, err error) {
	result, err := rpc.ProductClient.UploadImage(h.Context, &rpcproduct.UploadImageReq{
		FileName:  req.Name,
		ImageData: req.Image,
		Target:    req.Target,
	})
	if err != nil {
		return nil, err
	}
	return &product.UploadImageResp{
		Url: result.ImageUrl,
	}, nil
}
