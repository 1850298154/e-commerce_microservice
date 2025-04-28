package main

import (
	"fmt"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
	"time"

	"2501YTC/app/checkout/conf"
	"2501YTC/app/checkout/infra/rpc"
	"2501YTC/rpc_gen/kitex_gen/checkout/checkoutservice"

	// 导入 handler 包
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	opts := kitexInit()
	rpc.InitClient()

	svr := checkoutservice.NewServer(new(CheckoutServiceImpl), opts...)
	fmt.Println("checkout service start...")

	err := svr.Run()
	fmt.Println("checkout service start...")
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	// r, err := consul.NewConsulRegister(conf.GetConf().Registry.RegistryAddress[0])
	// if err != nil {
	//	 panic(err)
	// }

	// 服务注册与发现
	r, err := etcd.NewEtcdRegistry(conf.GetConf().Registry.RegistryAddress) // r should not be reused.
	if err != nil {
		log.Fatal(err)
	}
	opts = append(opts, server.WithRegistry(r))

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	// klog
	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Kitex.LogFileName,
			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		}),
		FlushInterval: time.Minute,
	}
	klog.SetOutput(asyncWriter)
	server.RegisterShutdownHook(func() {
		_ = asyncWriter.Sync()
	})
	return
}
