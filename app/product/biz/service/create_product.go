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

// Run 创建商品信息
func (s *CreateProductService) Run(req *product.CreateProductReq) (resp *product.CreateProductResp, err error) {
	// 检查商品名称是否为空
	if req.Name == "" {
		return nil, apiErr.ProductNameRequiredErr
	}

	// 检查商品价格是否合法
	if req.Price < 0 {
		return nil, apiErr.ProductPriceInvalidErr
	}

	// 初始化数据库查询对象
	q := model.NewProductQuery(s.ctx, mysql.DB)
	// 创建带缓存的查询对象
	cq := model.NewCachedProductQuery(q, redis.RedisClient)
	// 创建新商品
	p, err := cq.Create(model.Product{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  strings.Join(req.Categories, "+"), // 将分类数组用+号连接成字符串
	})
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	// 获取所有商品用于更新搜索引擎
	products, err := cq.GetAll()
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	// 将商品数据添加到Meilisearch搜索引擎
	err = meili.Add(products)
	if err != nil {
		return nil, apiErr.ConvertErr(err)
	}
	// 返回新创建的商品ID
	return &product.CreateProductResp{
		Id: uint32(p.ID),
	}, nil
}
