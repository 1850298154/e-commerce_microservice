package service

import (
	"context"
	"fmt"
	"strconv"

	"2501YTC/app/cart/biz/dal/redis"
	"2501YTC/app/cart/biz/model"
	"2501YTC/app/cart/errno"
	cart "2501YTC/rpc_gen/kitex_gen/cart"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
)

type GetCartService struct {
	ctx context.Context
} // NewGetCartService new GetCartService
func NewGetCartService(ctx context.Context) *GetCartService {
	return &GetCartService{ctx: ctx}
}

// Run create note info
func (s *GetCartService) Run(req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	// 从数据库中查找购物车列表
	cartService := model.GetCartService(redis.RedisClient)

	cartList, err := cartService.GetCartByUserId(s.ctx, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "%v", errno.GetCartErr(err))
		return nil, kerrors.NewBizStatusError(errno.GetCartErrCode, err.Error())
	}
	items := make([]*cart.CartItem, 0, len(cartList))
	// 将购物车列表转换为rpc返回的格式
	for _, v := range cartList {
		// 将字符串转换为 uint64
		u64, err := strconv.ParseUint(v.ProductId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID: %s", v.ProductId)
		}

		// 将 uint64 转换为 uint32
		u32 := uint32(u64)
		items = append(items, &cart.CartItem{
			ProductId: u32,
			Quantity:  v.Quantity,
		})
	}

	return &cart.GetCartResp{Cart: &cart.Cart{UserId: req.GetUserId(), Items: items}}, nil
}
