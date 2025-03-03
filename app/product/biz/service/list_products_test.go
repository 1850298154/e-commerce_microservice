package service

import (
	"context"
	"testing"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

func TestListProducts_Run(t *testing.T) {
	// 创建测试用例
	tests := []struct {
		name    string
		req     *product.ListProductsReq
		wantErr bool
	}{
		{
			name: "正常获取商品列表",
			req: &product.ListProductsReq{
				Page:     1,
				PageSize: 10,
			},
			wantErr: false,
		},
		{
			name: "页码为负数",
			req: &product.ListProductsReq{
				Page:     -1,
				PageSize: 10,
			},
			wantErr: true,
		},
		{
			name: "每页数量为负数",
			req: &product.ListProductsReq{
				Page:     1,
				PageSize: -1,
			},
			wantErr: true,
		},
		{
			name: "每页数量为0",
			req: &product.ListProductsReq{
				Page:     1,
				PageSize: 0,
			},
			wantErr: true,
		},
		{
			name: "按分类筛选",
			req: &product.ListProductsReq{
				Page:         1,
				PageSize:     10,
				CategoryName: "电子产品",
			},
			wantErr: false,
		},
		{
			name: "空分类名称",
			req: &product.ListProductsReq{
				Page:         1,
				PageSize:     10,
				CategoryName: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewListProductsService(ctx)

			resp, err := s.Run(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ListProducts.Run() 期望错误但是没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("ListProducts.Run() 错误 = %v", err)
				return
			}

			if resp == nil {
				t.Error("ListProducts.Run() 响应不应该为空")
				return
			}

			if resp.Products == nil {
				t.Error("ListProducts.Run() 返回的商品列表不应该为空")
			}
		})
	}
}
