package main

import (
	"context"
	"flag"
	"net"
	"time"

	"github.com/kitex-contrib/obs-opentelemetry/provider"

	"2501YTC/app/user/biz/dal"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"

	"github.com/joho/godotenv"
	consul "github.com/kitex-contrib/registry-consul"

	"2501YTC/app/user/conf"
	"2501YTC/rpc_gen/kitex_gen/user/userservice"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// 读取环境变量
	_ = godotenv.Load()

	// 初始化数据库服务
	dal.Init()

	// 处理命令行参数，运行不同的服务实例
	//port := flag.Int("port", 8082, "Service port")
	// weight := flag.Int("weight", 1, "Service weight")
	flag.Parse()

	// 初始化kitex服务
	opts := kitexInit()

	// 解析服务地址
	//addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *port))
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	// 链路追踪
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(conf.GetConf().Kitex.Service),
		// Support setting ExportEndpoint via environment variables: OTEL_EXPORTER_OTLP_ENDPOINT
		provider.WithExportEndpoint(conf.GetConf().OpenTelemetry.Endpoint),
		provider.WithInsecure(),
	)
	defer func(p provider.OtelProvider, ctx context.Context) {
		err := p.Shutdown(ctx)
		if err != nil {
			klog.Error(err.Error())
		}
	}(p, context.Background())

	svr := userservice.NewServer(new(UserServiceImpl), opts...)

	err = svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	// addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	// if err != nil {
	//	 panic(err)
	// }
	// opts = append(opts, server.WithServiceAddr(addr))

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	// 服务注册与发现
	r, err := consul.NewConsulRegister(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithRegistry(r))

	// 限流处理
	opts = append(opts, server.WithLimit(&limit.Option{
		MaxConnections: conf.GetConf().Kitex.MaxConnections, // MaxConnections: 最大连接数
		MaxQPS:         conf.GetConf().Kitex.MaxQPS,         // MaxQPS: 每秒最大请求数
	}))

	opts = append(opts, server.WithSuite(tracing.NewServerSuite()))

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
