package service

import (
	"context"
	"math"
	"strings"

	"2501YTC/app/product/biz/dal/mysql"
	"2501YTC/app/product/biz/model"
	"2501YTC/app/product/utils/apiErr"
	product "2501YTC/rpc_gen/kitex_gen/product"
)

type SearchProductsByNameService struct {
	ctx context.Context
} // NewSearchProductsByNameService new SearchProductsByNameService
func NewSearchProductsByNameService(ctx context.Context) *SearchProductsByNameService {
	return &SearchProductsByNameService{ctx: ctx}
}

// Run create note info
func (s *SearchProductsByNameService) Run(req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	// 检查搜索参数
	if req.Query == "" {
		return nil, apiErr.ProductNameRequiredErr
	}
	if req.Page <= 0 {
		return nil, apiErr.PageRequiredErr
	}
	if req.PageSize <= 0 {
		return nil, apiErr.PageSizeRequiredErr
	}
	// 初始化数据库查询对象
	q := model.NewProductQuery(s.ctx, mysql.DB)
	products, total, err := q.GetByName(req.Query, req.Page, req.PageSize)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	// 构建商品列表响应
	result := make([]*product.Product, 0)
	for _, p := range products {
		result = append(result, &product.Product{
			Id:         uint32(p.ID),
			Name:       p.Name,
			Price:      p.Price,
			Picture:    p.Picture,
			Categories: strings.Split(p.Categories, "+"), // 将分类字符串按+号分割为数组
		})
	}

	// 返回商品列表和总页数
	return &product.SearchProductsResp{
		Results: result,
		Num:     int64(math.Ceil(float64(total) / float64(req.PageSize))), // 计算总页数
	}, nil
}
