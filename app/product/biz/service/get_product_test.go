package service

import (
	"context"
	"testing"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

func TestGetProduct_Run(t *testing.T) {
	// 创建测试用例
	tests := []struct {
		name    string
		req     *product.GetProductReq
		wantErr bool
	}{
		{
			name: "正常获取商品",
			req: &product.GetProductReq{
				Id: 1,
			},
			wantErr: false,
		},
		{
			name: "获取不存在的商品",
			req: &product.GetProductReq{
				Id: 99999,
			},
			wantErr: true,
		},
		{
			name: "商品ID为0",
			req: &product.GetProductReq{
				Id: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewGetProductService(ctx)

			resp, err := s.Run(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetProduct.Run() 期望错误但是没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("GetProduct.Run() 错误 = %v", err)
				return
			}

			if resp == nil {
				t.Error("GetProduct.Run() 响应不应该为空")
				return
			}

			if resp.Product == nil {
				t.Error("GetProduct.Run() 返回的商品不应该为空")
			}
		})
	}
}
