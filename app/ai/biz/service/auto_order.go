package service

import (
	"context"
	ai "2501YTC/rpc_gen/kitex_gen/ai"
)

type AutoOrderService struct {
	ctx context.Context
} // NewAutoOrderService new AutoOrderService
func NewAutoOrderService(ctx context.Context) *AutoOrderService {
	return &AutoOrderService{ctx: ctx}
}

// Run create note info
func (s *AutoOrderService) Run(req *ai.AutoOrderReq) (resp *ai.AutoOrderResp, err error) {
	// Finish your business logic.

	return
}
