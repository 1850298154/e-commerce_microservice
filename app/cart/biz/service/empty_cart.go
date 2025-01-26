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
	// 清空购物车
	err = model.Cart.EmptyCart(model.Cart{}, s.ctx, redis.RedisClient, req.UserId)
	// err = model.Cart.EmptyCart(s.ctx, mysql.DB, req.UserId)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "清空购物车失败")
	}
	return &cart.EmptyCartResp{}, nil
}
