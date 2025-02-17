package service

import (
	"context"
	"strconv"

	"2501YTC/app/cart/biz/dal/redis"
	"2501YTC/app/cart/biz/model"
	"2501YTC/app/cart/errno"
	"2501YTC/app/cart/infra/rpc"

	cart "2501YTC/rpc_gen/kitex_gen/cart"
	"2501YTC/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
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
		klog.CtxErrorf(s.ctx, "%v", errno.ProductNotFoundErr(err))
		return nil, kerrors.NewBizStatusError(errno.ProductNotFoundErrCode, "商品未找到")
	}

	cartItem := &model.Cart{
		UserId:    req.UserId,
		ProductId: strconv.FormatUint(uint64(req.Item.ProductId), 10),
		Quantity:  req.Item.Quantity,
	}

	cartService := model.GetCartService(redis.RedisClient)
	err = cartService.AddItem(s.ctx, cartItem)
	if err != nil {
		klog.CtxErrorf(s.ctx, "%v", errno.GetCartErr(err))
		return nil, kerrors.NewBizStatusError(errno.AddToCartErrCode, "添加商品错误")
	}
	return &cart.AddItemResp{}, nil
}
