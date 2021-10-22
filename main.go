package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"

	"github.com/Junedayday/micro_web_service/gen/idl/demo"
	"github.com/Junedayday/micro_web_service/gen/idl/order"
	"github.com/Junedayday/micro_web_service/internal/config"
	"github.com/Junedayday/micro_web_service/internal/mysql"
	"github.com/Junedayday/micro_web_service/internal/server"
	"github.com/Junedayday/micro_web_service/internal/zlog"
)

var grpcGatewayTag = opentracing.Tag{Key: string(ext.Component), Value: "grpc-gateway"}

func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parentSpanContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err == nil || err == opentracing.ErrSpanContextNotFound {
			serverSpan := opentracing.GlobalTracer().StartSpan(
				"ServeHTTP",
				// this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
				ext.RPCServerOption(parentSpanContext),
				grpcGatewayTag,
			)
			r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))

			trace, ok := serverSpan.Context().(jaeger.SpanContext)
			if ok {
				w.Header().Set(jaeger.TraceContextHeaderName, fmt.Sprint(trace.TraceID()))
			}

			defer serverSpan.Finish()
		}
		h.ServeHTTP(w, r)
	})
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	// mux := runtime.NewServeMux(runtime.WithMetadata(annotator))
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		),
	}

	if err := demo.RegisterDemoServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", config.Viper.GetInt("server.grpc.port")), opts); err != nil {
		return errors.Wrap(err, "RegisterDemoServiceHandlerFromEndpoint error")
	} else if err := order.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", config.Viper.GetInt("server.grpc.port")), opts); err != nil {
		return errors.Wrap(err, "RegisterOrderServiceHandlerFromEndpoint error")
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Viper.GetInt("server.http.port")), tracingWrapper(mux))
}

func main() {
	var configFilePath = flag.String("c", "./", "config file path")
	flag.Parse()

	if err := config.Load(*configFilePath); err != nil {
		panic(err)
	}

	zlog.Init(config.Viper.GetString("zlog.path"))
	defer zlog.Sync()
	zlog.Sugar.Info("server is running")

	traceCfg := &jaegerconfig.Configuration{
		ServiceName: "MyService",
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LocalAgentHostPort: "127.0.0.1:6831",
			LogSpans:           true,
		},
	}
	tracer, closer, err := traceCfg.NewTracer(jaegerconfig.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// mysql初始化失败的话，不要继续运行程序
	if err := mysql.Init(config.Viper.GetString("mysql.user"),
		config.Viper.GetString("mysql.password"),
		config.Viper.GetString("mysql.ip"),
		config.Viper.GetInt("mysql.port"),
		config.Viper.GetString("mysql.dbname")); err != nil {
		zlog.Sugar.Fatalf("init mysql error %v", err)
	}

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Viper.GetInt("server.grpc.port")))
		if err != nil {
			panic(err)
		}

		s := grpc.NewServer(grpc.UnaryInterceptor(grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer()))))
		demo.RegisterDemoServiceServer(s, &server.Server{})
		order.RegisterOrderServiceServer(s, &server.Server{})

		if err = s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	if err := run(); err != nil {
		panic(err)
	}
}
