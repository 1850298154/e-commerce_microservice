package service

import (
	"context"
	"testing"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

func TestCreateProduct_Run(t *testing.T) {
	// 创建测试用例
	tests := []struct {
		name    string
		req     *product.CreateProductReq
		wantErr bool
	}{
		{
			name: "正常创建商品",
			req: &product.CreateProductReq{
				Name:       "测试商品",
				Price:      99.9,
				Picture:    "",
				Categories: []string{"电子产品", "手机"},
			},
			wantErr: false,
		},
		{
			name: "商品名称为空",
			req: &product.CreateProductReq{
				Name:       "",
				Price:      99.9,
				Picture:    "h",
				Categories: []string{"电子产品"},
			},
			wantErr: true,
		},
		{
			name: "商品价格为负",
			req: &product.CreateProductReq{
				Name:       "测试商品",
				Price:      -10,
				Picture:    "",
				Categories: []string{"电子产品"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewCreateProductService(ctx)

			resp, err := s.Run(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateProduct.Run() 期望错误但是没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateProduct.Run() 错误 = %v", err)
				return
			}

			if resp == nil {
				t.Error("CreateProduct.Run() 响应不应该为空")
				return
			}

			if resp.Id == 0 {
				t.Error("CreateProduct.Run() 创建的商品ID不应该为0")
			}
		})
	}
}
