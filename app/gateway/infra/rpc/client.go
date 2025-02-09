package rpc

import (
	"context"
	"sync"
	"time"

	"2501YTC/rpc_gen/kitex_gen/cart/cartservice"
	"2501YTC/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/fallback"
	"github.com/cloudwego/kitex/pkg/rpcinfo"

	gatewayutils "2501YTC/app/gateway/biz/utils"
	"2501YTC/rpc_gen/kitex_gen/auth/authservice"
	"2501YTC/rpc_gen/kitex_gen/product/productservice"
	"2501YTC/rpc_gen/kitex_gen/user/userservice"

	"2501YTC/app/gateway/conf"
	"2501YTC/common/clientsuite"
	"2501YTC/rpc_gen/kitex_gen/order/orderservice"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
)

const (
	serviceName        = "gateway"
	orderServiceName   = "order"
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
	OrderClient, err = orderservice.NewClient(orderServiceName, commonSuite)
	if err != nil {
		hlog.Fatal(err)
	}
}

func initUserClient() {
	var opts []client.Option

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
	// opts = append(opts, client.WithTracer())
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
