package service

import (
	"context"
	"strings"

	"2501YTC/app/product/biz/dal/redis"

	"2501YTC/app/product/biz/dal/mysql"
	"2501YTC/app/product/biz/model"
	"2501YTC/app/product/utils/apiErr"
	product "2501YTC/rpc_gen/kitex_gen/product"
)

type GetProductService struct {
	ctx context.Context
} // NewGetProductService new GetProductService
func NewGetProductService(ctx context.Context) *GetProductService {
	return &GetProductService{ctx: ctx}
}

// Run create note info
func (s *GetProductService) Run(req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	// Finish your business logic.
	if req.Id == 0 {
		return nil, apiErr.ProductIDRequiredErr
	}

	q := model.NewProductQuery(s.ctx, mysql.DB)
	cq := model.NewCachedProductQuery(q, redis.RedisClient)
	p, err := cq.GetById(int(req.Id))
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	return &product.GetProductResp{
		Product: &product.Product{
			Id:          uint32(p.ID),
			Picture:     p.Picture,
			Price:       p.Price,
			Description: p.Description,
			Name:        p.Name,
			Categories:  strings.Split(p.Categories, "+"),
		},
	}, nil
}
