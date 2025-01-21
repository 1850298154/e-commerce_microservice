package service

import (
	"context"
	payment "2501YTC/rpc_gen/kitex_gen/payment"
)

type CancelPaymentService struct {
	ctx context.Context
} // NewCancelPaymentService new CancelPaymentService
func NewCancelPaymentService(ctx context.Context) *CancelPaymentService {
	return &CancelPaymentService{ctx: ctx}
}

// Run create note info
func (s *CancelPaymentService) Run(req *payment.CancelPaymentReq) (resp *payment.CancelPaymentResp, err error) {
	// Finish your business logic.

	return
}
