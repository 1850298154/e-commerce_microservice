package service

import (
	"context"
	"strconv"

	"2501YTC/app/cart/biz/dal/redis"
	"2501YTC/app/cart/biz/model"
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
	// 检查商品是否存在
	productresp, err := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{Id: req.Item.ProductId})
	if err != nil {
		return nil, err
	}
	if productresp.Product == nil || productresp.Product.Id == 0 {
		return nil, kerrors.NewBizStatusError(10001, "商品不存在")
	}

	cartItem := &model.Cart{
		UserId:    req.UserId,
		ProductId: strconv.FormatUint(uint64(req.Item.ProductId), 10),
		Quantity:  req.Item.Quantity,
	}
	err = model.Cart.AddItem(model.Cart{}, s.ctx, redis.RedisClient, cartItem)

	//err = model.AddItem(s.ctx, mysql.DB, cartItem)

	if err != nil {
		return nil, kerrors.NewBizStatusError(10002, err.Error())
	}
	return &cart.AddItemResp{}, nil
}
