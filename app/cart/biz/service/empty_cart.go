package service

import (
	"context"

	"2501YTC/app/cart/biz/dal/redis"
	"2501YTC/app/cart/biz/model"
	cart "2501YTC/rpc_gen/kitex_gen/cart"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

type EmptyCartService struct {
	ctx context.Context
} // NewEmptyCartService new EmptyCartService
func NewEmptyCartService(ctx context.Context) *EmptyCartService {
	return &EmptyCartService{ctx: ctx}
}

// Run create note info
func (s *EmptyCartService) Run(req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	return &cart.EmptyCartResp{}, nil
	// 清空购物车
	cartService := model.GetCartService(redis.RedisClient)
	err = cartService.EmptyCart(s.ctx, req.UserId)
	if err != nil {
		return &cart.EmptyCartResp{}, kerrors.NewBizStatusError(50000, "清空购物车失败")
	}
	return &cart.EmptyCartResp{}, nil
}
