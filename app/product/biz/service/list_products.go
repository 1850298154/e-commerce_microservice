package service

import (
	"context"
	"math"
	"strings"

	"2501YTC/app/product/biz/dal/mysql"
	"2501YTC/app/product/biz/dal/redis"
	"2501YTC/app/product/biz/model"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

type ListProductsService struct {
	ctx context.Context
} // NewListProductsService new ListProductsService
func NewListProductsService(ctx context.Context) *ListProductsService {
	return &ListProductsService{ctx: ctx}
}

// Run create note info
func (s *ListProductsService) Run(req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	// Finish your business logic.
	q := model.NewProductQuery(s.ctx, mysql.DB)
	cq := model.NewCachedProductQuery(q, redis.RedisClient)
	products, total, err := cq.GetByCategory(req.CategoryName, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	result := make([]*product.Product, 0)
	for _, p := range products {
		result = append(result, &product.Product{
			Id:         uint32(p.ID),
			Name:       p.Name,
			Price:      p.Price,
			Picture:    p.Picture,
			Categories: strings.Split(p.Categories, "+"),
		})
	}
	return &product.ListProductsResp{
		Products: result,
		Num:      int64(math.Ceil(float64(total) / float64(req.PageSize))),
	}, nil
}
