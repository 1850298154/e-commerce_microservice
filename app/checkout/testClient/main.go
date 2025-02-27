package main

import (
	"context"
	"fmt"
	"log"

	"2501YTC/rpc_gen/kitex_gen/checkout"

	"2501YTC/app/checkout/conf"
	"2501YTC/rpc_gen/kitex_gen/checkout/checkoutservice"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
)

func main() {
	// 服务发现
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	c, err := checkoutservice.NewClient("checkout", client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("checkout client start...")
	resp, err := c.Checkout(ctx, &checkout.CheckoutReq{
		UserId:     0,
		Firstname:  "",
		Lastname:   "",
		Email:      "",
		Address:    nil,
		CreditCard: nil,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", resp)
}
