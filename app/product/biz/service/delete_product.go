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

// Run 删除商品信息
func (s *DeleteProductService) Run(req *product.DeleteProductReq) (resp *product.DeleteProductResp, err error) {
	// 检查商品ID是否有效
	if req.Id == 0 {
		return nil, apiErr.ProductIDRequiredErr
	}

	// 初始化数据库查询对象
	q := model.NewProductQuery(s.ctx, mysql.DB)
	// 创建带缓存的查询对象
	cq := model.NewCachedProductQuery(q, redis.RedisClient)

	// 根据ID获取商品信息
	p, err := cq.GetById(int(req.Id))
	if err != nil {
		// 商品不存在的情况
		return &product.DeleteProductResp{Success: false}, apiErr.ConvertErr(err)
	}

	// 异步删除商品图片
	if p.Picture != "" {
		img.DeleteObjectByUrlAsync(s.ctx, p.Picture)
	}

	// 从数据库中删除商品
	err = cq.Delete(int(p.ID))
	if err != nil {
		return &product.DeleteProductResp{Success: false}, apiErr.ConvertErr(err)
	}

	// 获取所有商品用于更新搜索引擎
	products, err := cq.GetAll()
	if err != nil {
		return &product.DeleteProductResp{Success: false}, apiErr.ConvertErr(err)
	}

	// 更新搜索引擎中的商品数据
	err = meili.Add(products)
	if err != nil {
		return &product.DeleteProductResp{Success: false}, apiErr.ConvertErr(err)
	}

	return &product.DeleteProductResp{Success: true}, nil
}
