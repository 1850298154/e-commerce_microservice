package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"2501YTC/app/checkout/infra/rpc"
	checkout "2501YTC/rpc_gen/kitex_gen/checkout"
	payment "2501YTC/rpc_gen/kitex_gen/payment"
)

func TestCheckout_Run(t *testing.T) {
	fmt.Println("TestCheckout_Run")
	// // 获取当前工作目录
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println("获取当前工作目录失败:", err)
	// 	return
	// }
	// fmt.Println("当前工作目录:", currentDir)

	// // 切换到新的工作目录
	// newDir := "D:\\zyt\\git_ln\\2501YTC\\app\\checkout"
	// err = os.Chdir(newDir)
	// if err != nil {
	// 	fmt.Println("切换工作目录失败:", err)
	// 	return
	// }

	// 获取当前文件的路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("获取当前文件路径失败")
		return
	}
	fmt.Println("当前文件路径:", filename)

	// 获取当前文件的目录
	currentDir := filepath.Dir(filename)
	fmt.Println("当前文件所在目录:", currentDir)

	// 计算上层路径（../..）
	upperDir := filepath.Join(currentDir, "..", "..")
	upperDir, err := filepath.Abs(upperDir) // 获取绝对路径
	if err != nil {
		fmt.Println("获取上层路径失败:", err)
		return
	}
	fmt.Println("上层路径:", upperDir)

	// 切换到新的工作目录
	err = os.Chdir(upperDir)
	if err != nil {
		fmt.Println("切换工作目录失败:", err)
		return
	}

	// 获取切换后的工作目录
	finalDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前工作目录失败:", err)
		return
	}
	fmt.Println("切换后的工作目录:", finalDir)

	// ===================================

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
