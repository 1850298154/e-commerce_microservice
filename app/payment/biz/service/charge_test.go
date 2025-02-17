package service

import (
	"context"
	"testing"

	"2501YTC/app/payment/biz/dal"
	payment "2501YTC/rpc_gen/kitex_gen/payment"
)

func TestCharge_Run(t *testing.T) {
	ctx := context.Background()
	s := NewChargeService(ctx)
	// init req and assert value
	// dsn := "root:root@tcp(127.0.0.1:3306)/payment?charset=utf8mb4&parseTime=True&loc=Local"
	dal.Init()
	req := &payment.ChargeReq{
		UserId:  1,
		OrderId: "1",
		Amount:  100,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          "4111111111111111",
			CreditCardCvv:             123,
			CreditCardExpirationMonth: 12,
			CreditCardExpirationYear:  2025,
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
