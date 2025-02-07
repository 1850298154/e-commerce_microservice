package service

import (
	"context"
	"testing"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

func TestDeleteProduct_Run(t *testing.T) {
	// 创建测试用例
	tests := []struct {
		name    string
		req     *product.DeleteProductReq
		wantErr bool
	}{
		{
			name: "正常删除商品",
			req: &product.DeleteProductReq{
				Id: 1,
			},
			wantErr: false,
		},
		{
			name: "删除不存在的商品",
			req: &product.DeleteProductReq{
				Id: 99999,
			},
			wantErr: true,
		},
		{
			name: "商品ID为0",
			req: &product.DeleteProductReq{
				Id: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewDeleteProductService(ctx)
			
			resp, err := s.Run(tt.req)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("DeleteProduct.Run() 期望错误但是没有错误")
				}
				return
			}
			
			if err != nil {
				t.Errorf("DeleteProduct.Run() 错误 = %v", err)
				return
			}
			
			if resp == nil {
				t.Error("DeleteProduct.Run() 响应不应该为空")
			}
		})
	}
}
