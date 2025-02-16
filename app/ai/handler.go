package main

import (
	ai "2501YTC/rpc_gen/kitex_gen/ai"
	"context"
	"2501YTC/app/ai/biz/service"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// QueryOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) QueryOrder(ctx context.Context, req *ai.OrderQueryReq) (resp *ai.OrderQueryResp, err error) {
	resp, err = service.NewQueryOrderService(ctx).Run(req)

	return resp, err
}

// AutoOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) AutoOrder(ctx context.Context, req *ai.AutoOrderReq) (resp *ai.AutoOrderResp, err error) {
	resp, err = service.NewAutoOrderService(ctx).Run(req)

	return resp, err
}
