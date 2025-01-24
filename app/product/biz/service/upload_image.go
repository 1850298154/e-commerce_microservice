package service

import (
	"bytes"
	"context"
	"os"

	"2501YTC/app/product/utils/apiErr"

	"2501YTC/app/product/utils/img"

	"github.com/cloudwego/kitex/pkg/klog"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

type UploadImageService struct {
	ctx context.Context
} // NewUploadImageService new UploadImageService
func NewUploadImageService(ctx context.Context) *UploadImageService {
	return &UploadImageService{ctx: ctx}
}

// Run create note info
func (s *UploadImageService) Run(req *product.UploadImageReq) (resp *product.UploadImageResp, err error) {
	// Finish your business logic.
	reader := bytes.NewReader(req.ImageData)

	size := int64(reader.Len())
	contentType := "image/jpeg"
	tempFile := "data/temp_compressed.jpg"
	err = img.ConvertAndCompressImage(s.ctx, reader, tempFile)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			klog.CtxErrorf(s.ctx, "failed to remove temp file: %s", err)
		}
	}(tempFile)

	compressedFile, err := os.Open(tempFile)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	defer func(compressedFile *os.File) {
		err := compressedFile.Close()
		if err != nil {
			klog.CtxErrorf(s.ctx, "failed to close compressed file: %s", err)
		}
	}(compressedFile)

	objectKey := img.GenerateObjectKey("image", ".jpeg")
	objectUrl, err := img.PutObject(objectKey, reader, size, contentType)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	return &product.UploadImageResp{
		ImageUrl: objectUrl,
	}, nil
}
