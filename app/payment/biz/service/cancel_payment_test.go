package service

import (
	"context"
	"testing"
	payment "2501YTC/rpc_gen/kitex_gen/payment"
)

func TestCancelPayment_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCancelPaymentService(ctx)
	// init req and assert value

	req := &payment.CancelPaymentReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
