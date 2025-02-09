package service

import (
	"context"
	"fmt"
	"testing"

	checkout "2501YTC/rpc_gen/kitex_gen/checkout"
	payment "2501YTC/rpc_gen/kitex_gen/payment"
)

func TestCheckout_Run(t *testing.T) {
	fmt.Println("TestCheckout_Run")
	ctx := context.Background()
	s := NewCheckoutService(ctx)
	// init req and assert value

	req := &checkout.CheckoutReq{
		UserId:    111,
		Firstname: "testFirstname",
		Lastname:  "testLastname",
		Email:     "testEmail",
		Address: &checkout.Address{
			StreetAddress: "testStreetAddress",
			City:          "testCity",
			State:         "testState",
			Country:       "testCountry",
			ZipCode:       "testZipCode",
		},
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          "testCreditCardNumber",
			CreditCardCvv:             3,
			CreditCardExpirationYear:  3,
			CreditCardExpirationMonth: 3,
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}
