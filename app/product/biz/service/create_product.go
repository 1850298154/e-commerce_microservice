package service

import (
	"context"
	"strings"

	"2501YTC/app/product/biz/dal/redis"

	"2501YTC/app/product/biz/dal/meili"
	"2501YTC/app/product/utils/apiErr"

	"2501YTC/app/product/biz/dal/mysql"
	"2501YTC/app/product/biz/model"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

type CreateProductService struct {
	ctx context.Context
} // NewCreateProductService new CreateProductService
func NewCreateProductService(ctx context.Context) *CreateProductService {
	return &CreateProductService{ctx: ctx}
}

// Run create note info
func (s *CreateProductService) Run(req *product.CreateProductReq) (resp *product.CreateProductResp, err error) {
	// Finish your business logic.
	q := model.NewProductQuery(s.ctx, mysql.DB)
	cq := model.NewCachedProductQuery(q, redis.RedisClient)
	p, err := cq.Create(model.Product{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  strings.Join(req.Categories, "+"),
	})
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	products, err := cq.GetAll()
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	err = meili.Add(products)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	return &product.CreateProductResp{
		Id: uint32(p.ID),
	}, nil
}
