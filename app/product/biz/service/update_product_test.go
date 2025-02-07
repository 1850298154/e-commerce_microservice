package service

import (
	"context"
	"testing"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

func TestUpdateProduct_Run(t *testing.T) {
	// 创建测试用例
	tests := []struct {
		name    string
		req     *product.UpdateProductReq
		wantErr bool
	}{
		{
			name: "正常更新商品",
			req: &product.UpdateProductReq{
				Id:          1,
				Name:        "iPhone 15",
				Price:       6999.00,
				Picture:     "http://example.com/iphone15.jpg",
				Description: "最新款iPhone",
				Categories:  []string{"手机", "电子产品"},
			},
			wantErr: false,
		},
		{
			name: "商品ID为0",
			req: &product.UpdateProductReq{
				Id:          0,
				Name:        "iPhone 15",
				Price:       6999.00,
				Picture:     "http://example.com/iphone15.jpg",
				Description: "最新款iPhone",
				Categories:  []string{"手机", "电子产品"},
			},
			wantErr: true,
		},
		{
			name: "商品名称为空",
			req: &product.UpdateProductReq{
				Id:          1,
				Name:        "",
				Price:       6999.00,
				Picture:     "http://example.com/iphone15.jpg",
				Description: "最新款iPhone",
				Categories:  []string{"手机", "电子产品"},
			},
			wantErr: true,
		},
		{
			name: "商品价格为负",
			req: &product.UpdateProductReq{
				Id:          1,
				Name:        "iPhone 15",
				Price:       -100,
				Picture:     "http://example.com/iphone15.jpg",
				Description: "最新款iPhone",
				Categories:  []string{"手机", "电子产品"},
			},
			wantErr: true,
		},
		{
			name: "商品图片为空",
			req: &product.UpdateProductReq{
				Id:          1,
				Name:        "iPhone 15",
				Price:       6999.00,
				Picture:     "",
				Description: "最新款iPhone",
				Categories:  []string{"手机", "电子产品"},
			},
			wantErr: true,
		},
		{
			name: "商品分类为空",
			req: &product.UpdateProductReq{
				Id:          1,
				Name:        "iPhone 15",
				Price:       6999.00,
				Picture:     "http://example.com/iphone15.jpg",
				Description: "最新款iPhone",
				Categories:  []string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUpdateProductService(context.Background())
			_, err := s.Run(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateProduct.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
