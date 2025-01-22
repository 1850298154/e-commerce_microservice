package service

import (
	"context"
	"strings"

	"2501YTC/app/product/utils/img"

	"2501YTC/app/product/biz/dal/redis"

	"2501YTC/app/product/biz/dal/meili"
	"2501YTC/app/product/biz/dal/mysql"
	"2501YTC/app/product/biz/model"
	"2501YTC/app/product/utils/apiErr"

	product "2501YTC/rpc_gen/kitex_gen/product"
)

type UpdateProductService struct {
	ctx context.Context
} // NewUpdateProductService new UpdateProductService
func NewUpdateProductService(ctx context.Context) *UpdateProductService {
	return &UpdateProductService{ctx: ctx}
}

// Run create note info
func (s *UpdateProductService) Run(req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
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
	if p.Picture != req.Picture {
		img.DeleteObjectByUrlAsync(s.ctx, p.Picture)
	}
	_, err = cq.Update(uint(req.Id), model.Product{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  strings.Join(req.Categories, "+"),
	})
	if err != nil {
		return &product.UpdateProductResp{Success: false}, apiErr.ConvertErr(err)
	}
	products, err := cq.GetAll()
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	err = meili.Add(products)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	return &product.UpdateProductResp{Success: true}, nil
}
