package service

import (
	ai "2501YTC/rpc_gen/kitex_gen/ai"
	"context"
)

type SearchforOrderService struct {
	ctx context.Context
} // NewSearchforOrderService new SearchforOrderService
func NewSearchforOrderService(ctx context.Context) *SearchforOrderService {
	return &SearchforOrderService{ctx: ctx}
}

// Run create note info
func (s *SearchforOrderService) Run(req *ai.SearchforOrderReq) (resp *ai.SearchforOrderResp, err error) {

	return
}
