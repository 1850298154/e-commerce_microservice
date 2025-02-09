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

// Run 上传图片服务
func (s *UploadImageService) Run(req *product.UploadImageReq) (resp *product.UploadImageResp, err error) {
	// 检查图片数据是否为空
	if len(req.ImageData) == 0 {
		return nil, apiErr.ImageDataRequiredErr
	}

	// 检查文件名是否为空
	if req.FileName == "" {
		return nil, apiErr.FileNameRequiredErr
	}

	// 创建图片数据读取器
	reader := bytes.NewReader(req.ImageData)

	// 获取图片大小
	size := int64(reader.Len())
	contentType := "image/jpeg"
	tempFile := "data/temp_compressed.jpg"

	// 转换并压缩图片
	err = img.ConvertAndCompressImage(s.ctx, reader, tempFile)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	// 延迟删除临时文件
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			klog.CtxErrorf(s.ctx, "删除临时文件失败: %s", err)
		}
	}(tempFile)

	// 打开压缩后的文件
	compressedFile, err := os.Open(tempFile)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	// 延迟关闭文件
	defer func(compressedFile *os.File) {
		err := compressedFile.Close()
		if err != nil {
			klog.CtxErrorf(s.ctx, "关闭压缩文件失败: %s", err)
		}
	}(compressedFile)

	// 生成对象存储的key并上传
	objectKey := img.GenerateObjectKey("image", ".jpeg")
	objectUrl, err := img.PutObject(objectKey, reader, size, contentType)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}

	// 返回上传成功的图片URL
	return &product.UploadImageResp{
		ImageUrl: objectUrl,
	}, nil
}
