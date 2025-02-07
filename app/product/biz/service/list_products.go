package service

import (
	"context"
	"math"
	"strings"

	"2501YTC/app/product/biz/dal/mysql"
	"2501YTC/app/product/biz/dal/redis"
	"2501YTC/app/product/biz/model"
	"2501YTC/app/product/utils/apiErr"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

type ListProductsService struct {
	ctx context.Context
} // NewListProductsService new ListProductsService
func NewListProductsService(ctx context.Context) *ListProductsService {
	return &ListProductsService{ctx: ctx}
}

// Run 获取商品列表信息
func (s *ListProductsService) Run(req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	// 检查分页参数
	if req.Page <= 0 {
		return nil, apiErr.PageRequiredErr
	}
	if req.PageSize <= 0 {
		return nil, apiErr.PageSizeRequiredErr
	}

	// 初始化数据库查询对象
	q := model.NewProductQuery(s.ctx, mysql.DB)
	// 创建带缓存的查询对象
	cq := model.NewCachedProductQuery(q, redis.RedisClient)
	
	// 根据分类获取商品列表
	products, total, err := cq.GetByCategory(req.CategoryName, req.Page, req.PageSize)
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
	return &product.ListProductsResp{
		Products: result,
		Num:      int64(math.Ceil(float64(total) / float64(req.PageSize))), // 计算总页数
	}, nil
}
