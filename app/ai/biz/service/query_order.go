package service

import (
	"context"
	ai "2501YTC/rpc_gen/kitex_gen/ai"
)

type QueryOrderService struct {
	ctx context.Context
} // NewQueryOrderService new QueryOrderService
func NewQueryOrderService(ctx context.Context) *QueryOrderService {
	return &QueryOrderService{ctx: ctx}
}

// Run create note info
func (s *QueryOrderService) Run(req *ai.OrderQueryReq) (resp *ai.OrderQueryResp, err error) {
	// Finish your business logic.

	return
}
