package service

import (
	payment "2501YTC/rpc_gen/kitex_gen/payment"
	"context"
	"testing"
)

func TestTimedCancelPayment_Run(t *testing.T) {
	ctx := context.Background()
	s := NewTimedCancelPaymentService(ctx)
	// init req and assert value

	req := &payment.TimedCancelPaymentReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
