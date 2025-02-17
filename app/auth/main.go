package main

import (
	"2501YTC/app/auth/biz/dal"
	"context"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/joho/godotenv"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"

	"2501YTC/app/auth/conf"
	"2501YTC/rpc_gen/kitex_gen/auth/authservice"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	_ = godotenv.Load()
	dal.Init()
	opts := kitexInit()

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

	svr := authservice.NewServer(new(AuthServiceImpl), opts...)

	if err := svr.Run(); err != nil {
		klog.Fatal(err)
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
		if err := asyncWriter.Sync(); err != nil {
			err.Error()
		}
	})

	return
}
