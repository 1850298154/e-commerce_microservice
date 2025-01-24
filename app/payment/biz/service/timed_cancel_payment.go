package service

import (
	payment "2501YTC/rpc_gen/kitex_gen/payment"
	"context"
)

type TimedCancelPaymentService struct {
	ctx context.Context
} // NewTimedCancelPaymentService new TimedCancelPaymentService
func NewTimedCancelPaymentService(ctx context.Context) *TimedCancelPaymentService {
	return &TimedCancelPaymentService{ctx: ctx}
}

// Run create note info
func (s *TimedCancelPaymentService) Run(req *payment.TimedCancelPaymentReq) (resp *payment.TimedCancelPaymentResp, err error) {
	// Finish your business logic.

	return
}
