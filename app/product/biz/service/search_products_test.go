package service

import (
	"context"
	"fmt"
	"testing"

	"2501YTC/app/product/biz/dal"

	"github.com/joho/godotenv"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

func TestSearchProducts_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	// 创建测试用例
	tests := []struct {
		name    string
		req     *product.SearchProductsReq
		wantErr bool
	}{
		{
			name: "正常搜索商品",
			req: &product.SearchProductsReq{
				Query:    "电脑",
				Page:     1,
				PageSize: 10,
			},
			wantErr: false,
		},
		{
			name: "空搜索关键词",
			req: &product.SearchProductsReq{
				Query:    "",
				Page:     1,
				PageSize: 10,
			},
			wantErr: true,
		},
		{
			name: "页码为负数",
			req: &product.SearchProductsReq{
				Query:    "手机",
				Page:     -1,
				PageSize: 10,
			},
			wantErr: true,
		},
		{
			name: "每页数量为负数",
			req: &product.SearchProductsReq{
				Query:    "手机",
				Page:     1,
				PageSize: -1,
			},
			wantErr: true,
		},
		{
			name: "每页数量为0",
			req: &product.SearchProductsReq{
				Query:    "手机",
				Page:     1,
				PageSize: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSearchProductsService(context.Background())
			resp, err := s.Run(tt.req)
			fmt.Println(resp.Results)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchProducts.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Error("SearchProducts.Run() 期望返回结果,但是得到了nil")
			}
		})
	}
}
