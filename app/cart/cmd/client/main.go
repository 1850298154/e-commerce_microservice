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
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

func main() {
	// 服务发现
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName("cart_client"),
		provider.WithExportEndpoint("localhost:4317"),
		provider.WithInsecure(),
	)
	defer func() {
		_ = p.Shutdown(ctx)
	}()
	// 链路追踪
	c, err := cartservice.NewClient("cart", client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()))
	// resp, err := c.AddItem(ctx, &cart.AddItemReq{UserId: 1, Item: &cart.CartItem{ProductId: 1, Quantity: 1}})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%+v", resp)

	for {
		getcartresp, err := c.GetCart(ctx, &cart.GetCartReq{UserId: 1})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("获取购物车成功：%+v", getcartresp)

		emptycartresp, err := c.EmptyCart(ctx, &cart.EmptyCartReq{UserId: 1})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("清空购物车成功：%+v", emptycartresp)
		<-time.After(2 * time.Second)
	}
}
