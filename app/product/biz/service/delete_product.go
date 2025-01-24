package service

import (
	"context"

	"2501YTC/app/product/utils/img"

	"2501YTC/app/product/biz/dal/redis"

	"2501YTC/app/product/biz/dal/meili"
	"2501YTC/app/product/biz/dal/mysql"
	"2501YTC/app/product/biz/model"
	"2501YTC/app/product/utils/apiErr"
	product "2501YTC/rpc_gen/kitex_gen/product"
)

type DeleteProductService struct {
	ctx context.Context
} // NewDeleteProductService new DeleteProductService
func NewDeleteProductService(ctx context.Context) *DeleteProductService {
	return &DeleteProductService{ctx: ctx}
}

// Run create note info
func (s *DeleteProductService) Run(req *product.DeleteProductReq) (resp *product.DeleteProductResp, err error) {
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
	img.DeleteObjectByUrlAsync(s.ctx, p.Picture)
	err = cq.Delete(int(p.ID))
	if err != nil {
		return &product.DeleteProductResp{Success: false}, apiErr.ConvertErr(err)
	}
	products, err := cq.GetAll()
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	err = meili.Add(products)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	return &product.DeleteProductResp{Success: true}, nil
}
