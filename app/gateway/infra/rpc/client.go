package rpc

import(
	"sync"

	"2501YTC/rpc_gen/kitex_gen/order/orderservice"
	"2501YTC/app/gateway/conf"
	"2501YTC/common/clientsuite"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const(
	serviceName = "gateway"
)

var(
	OrderClient orderservice.Client

	once sync.Once
	err error
	registryAddr string
	commonSuite client.Option
)

func InitClient() {
	once.Do(func() {
		registryAddr = conf.GetConf().Hertz.RegistryAddr
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr: registryAddr,
			CurrentServiceName: serviceName,
		})
		initOrderClient()
	})
}

func initOrderClient() {
	OrderClient, err = orderservice.NewClient("order", commonSuite)
	if err != nil {
		hlog.Fatal(err)
	}
}