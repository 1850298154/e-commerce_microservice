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

// Run 更新商品信息
func (s *UpdateProductService) Run(req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	// 检查商品ID是否有效
	if req.Id == 0 {
		return nil, apiErr.ProductIDRequiredErr
	}

	// 检查商品名称是否为空
	if req.Name == "" {
		return nil, apiErr.ProductNameRequiredErr
	}

	// 检查商品价格是否合法
	if req.Price < 0 {
		return nil, apiErr.ProductPriceInvalidErr
	}

	// 创建数据库查询对象
	q := model.NewProductQuery(s.ctx, mysql.DB)
	// 创建带缓存的查询对象
	cq := model.NewCachedProductQuery(q, redis.RedisClient)

	// 获取原商品信息
	p, err := cq.GetById(int(req.Id))
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}

	// 如果商品图片发生变化,异步删除旧图片
	if p.Picture != req.Picture {
		img.DeleteObjectByUrlAsync(s.ctx, p.Picture)
	}

	// 更新商品信息
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

	// 获取所有商品用于更新搜索引擎
	products, err := cq.GetAll()
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}

	// 更新搜索引擎数据
	err = meili.Add(products)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}

	return &product.UpdateProductResp{Success: true}, nil
}
