package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	"2501YTC/app/checkout/infra/rpc"
	checkout "2501YTC/rpc_gen/kitex_gen/checkout"
	payment "2501YTC/rpc_gen/kitex_gen/payment"
)

func TestCheckout_Run(t *testing.T) {
	fmt.Println("TestCheckout_Run")
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前工作目录失败:", err)
		return
	}
	fmt.Println("当前工作目录:", currentDir)

	// 切换到新的工作目录
	newDir := "D:\\zyt\\git_ln\\2501YTC\\app\\checkout"
	err = os.Chdir(newDir)
	if err != nil {
		fmt.Println("切换工作目录失败:", err)
		return
	}

	ctx := context.Background()
	s := NewCheckoutService(ctx)
	rpc.InitClient()

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
