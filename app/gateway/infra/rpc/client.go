package rpc

import (
	"context"
	"fmt"
	"sync"

	gatewayutils "2501YTC/app/gateway/biz/utils"
	"2501YTC/rpc_gen/kitex_gen/auth/authservice"
	"2501YTC/rpc_gen/kitex_gen/user/userservice"

	"2501YTC/app/gateway/conf"
	"2501YTC/common/clientsuite"
	"2501YTC/rpc_gen/kitex_gen/order/orderservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	// "go.opentelemetry.io/otel"
)

const (
	serviceName      = "gateway"
	orderServiceName = "order"
	OrderClientName  = "orderClient"
)

var (
	OrderClient  orderservice.Client
	UserClient   userservice.Client
	AuthClient   authservice.Client
	once         sync.Once
	err          error
	registryAddr string
	commonSuite  client.Option
)

func InitClient() {
	once.Do(func() {
		registryAddr = conf.GetConf().Hertz.RegistryAddr
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: serviceName,
		})
		initOrderClient()
		initUserClient()
		initAuthClient()
	})
}

func initOrderClient() {
	// TODO Opentelemetry
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(OrderClientName),
		// Support setting ExportEndpoint via environment variables: OTEL_EXPORTER_OTLP_ENDPOINT
		provider.WithExportEndpoint(":4317"),
		provider.WithInsecure(),
	)
	defer func() {
		_ = p.Shutdown(context.Background())
	}()

	// TODO 负载均衡、熔断
	var opts []client.Option
	// 熔断器配置
	// build a new CBSuite with default config CBConfig{Enable: true, ErrRate: 0.5, MinSample: 200}
	cbs := circuitbreak.NewCBSuite(func(ri rpcinfo.RPCInfo) string {
		return circuitbreak.RPCInfo2Key(ri)
	})
	// customize the circuit breaker config for the service
	cbs.UpdateServiceCBConfig(fmt.Sprintf("%s/%s/%s", serviceName, orderServiceName, "ListOrder"), circuitbreak.CBConfig{
		Enable:    true,
		ErrRate:   0.7, // requests will be blocked if error rate >= 30%
		MinSample: 400, // this config takes effect if sampled requests are more than `MinSample`
	})
	// 加入负载均衡、熔断器
	opts = append(opts, commonSuite, client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()), client.WithCircuitBreaker(cbs))
	// 加入tracing
	opts = append(opts, client.WithSuite(tracing.NewClientSuite()))
	// 加入rpcinfo
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: OrderClientName}))

	OrderClient, err = orderservice.NewClient(orderServiceName, opts...)
	gatewayutils.MustHandleError(err)
}

func initUserClient() {
	UserClient, err = userservice.NewClient("user", commonSuite)
	gatewayutils.MustHandleError(err)
}

func initAuthClient() {
	AuthClient, err = authservice.NewClient("auth", commonSuite)
	gatewayutils.MustHandleError(err)
}
