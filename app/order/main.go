package main

import (
	"context"
	"net"
	"time"

	"2501YTC/app/order/biz/dal"
	"2501YTC/app/order/biz/dal/mq"
	"2501YTC/app/order/biz/dal/mysql"
	"2501YTC/app/order/conf"
	"2501YTC/common/healthcheck"
	"2501YTC/rpc_gen/kitex_gen/order/orderservice"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
)

var consumer *mq.Consumer

func main() {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc

	// TODO opentelemetry
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(conf.GetConf().Kitex.Service),
		provider.WithExportEndpoint(conf.GetConf().OpenTelemetry.Endpoint),
		provider.WithInsecure(),
	)
	defer func() {
		_ = p.Shutdown(context.Background())
	}()

	// 健康检查
	healthcheck.StartHealthCheck(conf.GetConf().HealthCheck.Addr, conf.GetConf().Kitex.Service)
	klog.Infof("Health check server started on port %s", conf.GetConf().HealthCheck.Addr)

	// 初始化MySQL和RabbitMQ
	dal.Init()

	// 初始化Kitex
	opts := kitexInit()

	startProducer()
	defer mq.ProducerInstance.Stop()
	startConsumer(mysql.DB)
	defer consumer.Stop()

	svr := orderservice.NewServer(new(OrderServiceImpl), opts...)
	err := svr.Run()
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

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	// 服务注册
	r, err := consul.NewConsulRegister(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithRegistry(r))
	// TODO 限流
	opts = append(opts, server.WithLimit(&limit.Option{
		MaxConnections: conf.GetConf().Kitex.MaxConnections, // MaxConnections: 最大连接数
		MaxQPS:         conf.GetConf().Kitex.MaxQPS,         // MaxQPS: 每秒最大请求数
	}))
	// TODO tracing
	opts = append(opts, server.WithSuite(tracing.NewServerSuite()))
	// rpcinfo
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GetConf().Kitex.Service}))

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
		consumer.Stop()
		mq.ProducerInstance.Stop()
		_ = asyncWriter.Sync()
	})
	return
}

func startProducer() {
	klog.Info("Producer starting...")
	var err error
	mq.ProducerInstance, err = mq.NewProducer(conf.GetConf().RabbitMQ.OrderTimeout)
	if err != nil {
		klog.Fatalf("NewProducer failed, err: %v", err)
		panic(err)
	}
}

func startConsumer(db *gorm.DB) {
	klog.Info("Consumer starting...")
	var err error
	consumer, err = mq.NewConsumer(db)
	if err != nil {
		klog.Fatalf("NewConsumer failed, err: %v", err)
		panic(err)
	}
	go func() {
		_ = consumer.Consume()
	}()
}
