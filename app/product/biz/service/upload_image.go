package service

import (
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

// Run 上传图片服务
func (s *UploadImageService) Run(req *product.UploadImageReq) (resp *product.UploadImageResp, err error) {
	// 请求验证
	if len(req.ImageData) == 0 {
		return nil, apiErr.ImageDataRequiredErr
	}
	if req.FileName == "" {
		return nil, apiErr.FileNameRequiredErr
	}

	// 压缩图片处理
	tempFile := "data/temp_compressed.jpg"
	if err := img.ConvertAndCompressImage(s.ctx, req.ImageData, req.FileName, tempFile); err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	defer func() {
		if err := os.Remove(tempFile); err != nil {
			klog.CtxErrorf(s.ctx, "删除临时文件失败: %s", err)
		}
	}()
	// 打开压缩后的图片文件
	compressedFile, err := os.Open(tempFile)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	defer func() {
		if err := compressedFile.Close(); err != nil {
			klog.CtxErrorf(s.ctx, "关闭压缩文件失败: %s", err)
		}
	}()
	// 计算图片大小
	info, err := os.Stat(tempFile)
	if err != nil {
		return nil, err
	}
	// 上传图片到对象存储
	objectKey := img.GenerateObjectKey("image", ".jpeg")
	objectUrl, err := img.PutObject(objectKey, compressedFile, info.Size(), "image/jpeg")
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	// 返回上传成功的图片URL
	return &product.UploadImageResp{
		ImageUrl: objectUrl,
	}, nil
}
