package rpc

import (
	"2501YTC/app/gateway/conf"
	"2501YTC/common/clientsuite"
	"2501YTC/rpc_gen/kitex_gen/auth/authservice"
	"2501YTC/rpc_gen/kitex_gen/cart/cartservice"
	"2501YTC/rpc_gen/kitex_gen/order/orderservice"
	"2501YTC/rpc_gen/kitex_gen/product/productservice"
	"2501YTC/rpc_gen/kitex_gen/user"
	"2501YTC/rpc_gen/kitex_gen/user/userservice"
	"context"
	"fmt"
	"sync"
	"time"

	gatewayutils "2501YTC/app/gateway/biz/utils"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/fallback"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	// "go.opentelemetry.io/otel"
)

const (
	serviceName        = "gateway"
	orderServiceName   = "order"
	orderClientName    = "orderClient"
	userServiceName    = "user"
	authServiceName    = "auth"
	cartServiceName    = "cart"
	productServiceName = "product"
)

var (
	OrderClient   orderservice.Client
	UserClient    userservice.Client
	AuthClient    authservice.Client
	CartClient    cartservice.Client
	ProductClient productservice.Client
	once          sync.Once
	err           error
	registryAddr  string
	commonSuite   client.Option
)

func GenServiceCBKeyFunc(ri rpcinfo.RPCInfo) string {
	// circuitbreak.RPCInfo2Key returns "$fromServiceName/$toServiceName/$method"
	return circuitbreak.RPCInfo2Key(ri)
}

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
		initCartClient()
		initProductClient()
	})
}

func initOrderClient() {
	// TODO 负载均衡、熔断
	var opts []client.Option
	// 链路追踪
	opts = append(opts, client.WithSuite(tracing.NewClientSuite()))
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
	// 加入rpcinfo
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: orderClientName}))

	OrderClient, err = orderservice.NewClient(orderServiceName, opts...)
	gatewayutils.MustHandleError(err)
}

func initUserClient() {
	var opts []client.Option

	// 链路追踪
	opts = append(opts, client.WithSuite(tracing.NewClientSuite()))

	// 熔断机制
	cbs := circuitbreak.NewCBSuite(GenServiceCBKeyFunc)
	cbs.UpdateServiceCBConfig("gateway/user/Register", circuitbreak.CBConfig{Enable: true, ErrRate: 0.5, MinSample: 100000})
	cbs.UpdateServiceCBConfig("gateway/user/Login", circuitbreak.CBConfig{Enable: true, ErrRate: 0.5, MinSample: 100000})
	cbs.UpdateServiceCBConfig("gateway/user/GetUserInfo", circuitbreak.CBConfig{Enable: true, ErrRate: 0.5, MinSample: 100000})
	cbs.UpdateServiceCBConfig("gateway/user/UpdateUser", circuitbreak.CBConfig{Enable: true, ErrRate: 0.5, MinSample: 100000})

	registerFallback := func(ctx context.Context, req, resp any, err error) (fbResp any, fbErr error) {
		return &user.RegisterResp{UserId: 0}, nil
	}

	loginFallback := func(ctx context.Context, req, resp any, err error) (fbResp any, fbErr error) {
		return &user.LoginResp{
			UserId: 0,
			Role:   1,
		}, nil
	}

	getUserInfoFallback := func(ctx context.Context, req, resp any, err error) (fbResp any, fbErr error) {
		return &user.GetUserInfoResp{
			UserId:    0,
			Email:     "user@example.com",
			CreatedAt: time.Now().String(),
			UpdatedAt: time.Now().String(),
		}, nil
	}

	updateUserFallback := func(ctx context.Context, req, resp any, err error) (fbResp any, fbErr error) {
		return &user.UpdateUserResp{
			Success: false,
		}, nil
	}

	fallbackFunc := fallback.UnwrapHelper(func(ctx context.Context, req, resp any, err error) (fbResp any, fbErr error) {
		if err == nil {
			return resp, nil
		}
		methodName := rpcinfo.GetRPCInfo(ctx).To().Method()
		switch methodName {
		case "Register":
			return registerFallback(ctx, req, resp, err)
		case "Login":
			return loginFallback(ctx, req, resp, err)
		case "GetUserInfo":
			return getUserInfoFallback(ctx, req, resp, err)
		case "UpdateUser":
			return updateUserFallback(ctx, req, resp, err)
		default:
			return resp, err
		}
	})

	fallbackPolicy := fallback.NewFallbackPolicy(fallbackFunc)

	opts = append(opts, commonSuite, client.WithCircuitBreaker(cbs))
	opts = append(opts, client.WithFallback(fallbackPolicy))

	UserClient, err = userservice.NewClient(userServiceName, opts...)
	gatewayutils.MustHandleError(err)
}

func initAuthClient() {
	AuthClient, err = authservice.NewClient(authServiceName, commonSuite)
	gatewayutils.MustHandleError(err)
}

func initCartClient() {
	CartClient, err = cartservice.NewClient(cartServiceName, commonSuite)
	gatewayutils.MustHandleError(err)
}

func initProductClient() {
	ProductClient, err = productservice.NewClient(productServiceName, commonSuite)
	gatewayutils.MustHandleError(err)
}
