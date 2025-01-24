package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"2501YTC/app/cart/conf"
	"2501YTC/rpc_gen/kitex_gen/cart"
	"2501YTC/rpc_gen/kitex_gen/cart/cartservice"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
)

func main() {
	// 服务发现
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		log.Fatal(err)
	}
	c, err := cartservice.NewClient("cart", client.WithResolver(r), client.WithRPCTimeout(time.Second*3))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	resp, err := c.AddItem(ctx, &cart.AddItemReq{UserId: 1, Item: &cart.CartItem{ProductId: 2, Quantity: 1}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", resp)
}
