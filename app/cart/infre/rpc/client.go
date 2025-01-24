package rpc

import (
	"2501YTC/app/cart/conf"
	"2501YTC/rpc_gen/kitex_gen/product/productcatalogservice"
	"sync"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
)

var (
	ProductClient productcatalogservice.Client
	once          sync.Once
)

func InitClient() {
	once.Do(func() {
		InitProductClient()
	})
}

// 初始化RPC客户端
func InitProductClient() {
	// 声明客户端
	var opts []client.Option

	// 获取product服务
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		panic(err)
	}
	// 添加resolver
	opts = append(opts, client.WithResolver(r))

	// 创建product客户端
	ProductClient, err = productcatalogservice.NewClient("product", opts...)
}
