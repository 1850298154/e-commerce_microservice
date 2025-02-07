package rpc

import (
	"sync"

	gatewayutils "2501YTC/app/gateway/biz/utils"
	"2501YTC/rpc_gen/kitex_gen/auth/authservice"
	"2501YTC/rpc_gen/kitex_gen/user/userservice"

	"2501YTC/app/gateway/conf"
	"2501YTC/common/clientsuite"
	"2501YTC/rpc_gen/kitex_gen/order/orderservice"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
)

const (
	serviceName      = "gateway"
	orderServiceName = "order"
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
	OrderClient, err = orderservice.NewClient(orderServiceName, commonSuite)
	if err != nil {
		hlog.Fatal(err)
	}
}

func initUserClient() {
	UserClient, err = userservice.NewClient("user", commonSuite)
	gatewayutils.MustHandleError(err)
}

func initAuthClient() {
	AuthClient, err = authservice.NewClient("auth", commonSuite)
	gatewayutils.MustHandleError(err)
}
