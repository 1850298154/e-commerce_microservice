package main

import (
	"2501YTC/app/payment/biz/service"
	payment "2501YTC/rpc_gen/kitex_gen/payment"
	"context"
)

// PaymentServiceImpl implements the last service interface defined in the IDL.
type PaymentServiceImpl struct{}

// Charge implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) Charge(ctx context.Context, req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {
	resp, err = service.NewChargeService(ctx).Run(req)

	return resp, err
}

// CancelPayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) CancelPayment(ctx context.Context, req *payment.CancelPaymentReq) (resp *payment.CancelPaymentResp, err error) {
	resp, err = service.NewCancelPaymentService(ctx).Run(req)

	return resp, err
}

// TimedCancelPayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) TimedCancelPayment(ctx context.Context, req *payment.TimedCancelPaymentReq) (resp *payment.TimedCancelPaymentResp, err error) {
	resp, err = service.NewTimedCancelPaymentService(ctx).Run(req)

	return resp, err
}
