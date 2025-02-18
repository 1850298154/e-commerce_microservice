package rpc

import (
	"context"
	"fmt"
	"sync"

	"2501YTC/rpc_gen/kitex_gen/cart/cartservice"

	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"

	"2501YTC/common/clientsuite"
	"2501YTC/rpc_gen/kitex_gen/product/productservice"

	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/rpcinfo"

	gatewayutils "2501YTC/app/gateway/biz/utils"
	"2501YTC/rpc_gen/kitex_gen/order/orderservice"

	"github.com/cloudwego/kitex/client"
)

const (
	serviceName        = "ai"
	orderServiceName   = "order"
	orderClientName    = "orderClient"
	productServiceName = "product"
	productClientName  = "productClient"
	cartServiceName    = "cart"
	cartClientName     = "cartClient"
)

var (
	OrderClient   orderservice.Client
	ProductClient productservice.Client
	CartClient    cartservice.Client
	once          sync.Once
	err           error
	commonSuite   client.Option
)

func GenServiceCBKeyFunc(ri rpcinfo.RPCInfo) string {
	// circuitbreak.RPCInfo2Key returns "$fromServiceName/$toServiceName/$method"
	return circuitbreak.RPCInfo2Key(ri)
}

func InitClient() {
	once.Do(func() {
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       "localhost:8500",
			CurrentServiceName: serviceName,
		})
		initOrderClient()
		initProductClient()
		initCartClient()
	})
}

func initCartClient() {
	CartClient, err = cartservice.NewClient(cartServiceName, commonSuite)
	gatewayutils.MustHandleError(err)
}

func initOrderClient() {
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(orderClientName),
		provider.WithExportEndpoint(":4317"),
		provider.WithInsecure(),
	)
	defer func() {
		_ = p.Shutdown(context.Background())
	}()

	var opts []client.Option
	// 熔断器配置
	cbs := circuitbreak.NewCBSuite(GenServiceCBKeyFunc)
	cbs.UpdateServiceCBConfig(fmt.Sprintf("%s/%s/%s", serviceName, orderServiceName, "ListOrder"), circuitbreak.CBConfig{
		Enable:    true,
		ErrRate:   0.7,
		MinSample: 400,
	})
	// 加入负载均衡、熔断器
	opts = append(opts, commonSuite, client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()), client.WithCircuitBreaker(cbs))
	// 加入tracing
	opts = append(opts, client.WithSuite(tracing.NewClientSuite()))
	// 加入rpcinfo
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: orderClientName}))

	OrderClient, err = orderservice.NewClient(orderServiceName, opts...)
	gatewayutils.MustHandleError(err)
}

func initProductClient() {
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(productClientName),
		provider.WithExportEndpoint(":4317"),
		provider.WithInsecure(),
	)
	defer func() {
		_ = p.Shutdown(context.Background())
	}()

	var opts []client.Option
	// 熔断器配置
	cbs := circuitbreak.NewCBSuite(GenServiceCBKeyFunc)
	cbs.UpdateServiceCBConfig(fmt.Sprintf("%s/%s/%s", serviceName, productServiceName, "ListProduct"), circuitbreak.CBConfig{
		Enable:    true,
		ErrRate:   0.7,
		MinSample: 400,
	})
	// 加入负载均衡、熔断器
	opts = append(opts, commonSuite, client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()), client.WithCircuitBreaker(cbs))
	// 加入tracing
	opts = append(opts, client.WithSuite(tracing.NewClientSuite()))
	// 加入rpcinfo
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: productClientName}))
	ProductClient, err = productservice.NewClient(productServiceName, opts...)
	gatewayutils.MustHandleError(err)
}
