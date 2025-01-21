package service

import (
	"context"
	"fmt"

	"2501YTC/app/cart/infre/rpc"
	cart "2501YTC/rpc_gen/kitex_gen/cart"
	"2501YTC/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

type AddItemService struct {
	ctx context.Context
} // NewAddItemService new AddItemService
func NewAddItemService(ctx context.Context) *AddItemService {
	return &AddItemService{ctx: ctx}
}

// Run create note info
func (s *AddItemService) Run(req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	// 1. 检查商品是否存在
	// RPC调用获取商品请求
	productresp, err := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{Id: req.Item.ProductId})
	if err != nil {
		return nil, err
	}
	if productresp.Product == nil || productresp.Product.Id == 0 {
		return nil, kerrors.NewBizStatusError(10001, "商品不存在")
	}
	// 2. 检查商品是否已经在购物车中

	// 3. 检查商品是否已经购买

	// 4. 检查商品是否已经下架

	// 5. 检查商品是否已经删除

	// 6. 检查商品是否已经过期

	// 7. 检查商品是否已经达到最大购买数量

	// 8. 检查商品是否已经达到最大购买金额

	fmt.Printf("req: %v\n", req)
	return &cart.AddItemResp{}, nil
}
