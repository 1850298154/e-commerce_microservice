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
	// fmt.Println("GetProductService.Run")
	// // 构建并返回商品信息响应
	// return &product.GetProductResp{
	// 	Product: &product.Product{
	// 		Id:          req.Id,
	// 		Picture:     "Picture",
	// 		Price:       10,
	// 		Description: "Description",
	// 		Name:        "Name",
	// 		Categories:  []string{"category1", "category2", "category3"},
	// 	},
	// }, nil

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
		return nil, apiErr.ConvertErr(err)
	}

	// 构建并返回商品信息响应
	return &product.GetProductResp{
		Product: &product.Product{
			Id:          uint32(p.ID),
			Picture:     p.Picture,
			Price:       p.Price,
			Description: p.Description,
			Name:        p.Name,
			Categories:  strings.Split(p.Categories, "+"), // 将分类字符串按+号分割为数组
		},
	}, nil
}
