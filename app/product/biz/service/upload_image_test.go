package service

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"2501YTC/app/product/biz/dal"

	"github.com/joho/godotenv"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

func TestUploadImage_Run(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(filename), "../../")
	_ = godotenv.Load(basePath + "/.env")
	dal.Init()
	// 创建测试用例
	tests := []struct {
		name    string
		req     *product.UploadImageReq
		wantErr bool
	}{
		{
			name: "正常上传图片",
			req: &product.UploadImageReq{
				ImageData: []byte("test image data"),
				FileName:  "test.jpg",
			},
			wantErr: false,
		},
		{
			name: "空图片数据",
			req: &product.UploadImageReq{
				ImageData: []byte{},
				FileName:  "test.jpg",
			},
			wantErr: true,
		},
		{
			name: "不支持的图片类型",
			req: &product.UploadImageReq{
				ImageData: []byte("test image data"),
				FileName:  "test.invalid",
			},
			wantErr: true,
		},
		{
			name: "文件名为空",
			req: &product.UploadImageReq{
				ImageData: []byte("test image data"),
				FileName:  "",
			},
			wantErr: true,
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUploadImageService(context.Background())
			resp, err := s.Run(tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UploadImage.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && (resp == nil || resp.ImageUrl == "") {
				t.Error("UploadImage.Run() 返回的URL为空")
			}
		})
	}
}
