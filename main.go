package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/Junedayday/micro_web_service/gen/idl/demo"
	"github.com/Junedayday/micro_web_service/internal/config"
	"github.com/Junedayday/micro_web_service/internal/server"
	"github.com/Junedayday/micro_web_service/internal/zlog"
)

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := demo.RegisterDemoServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", config.Viper.GetInt("server.grpc.port")), opts); err != nil {
		return errors.Wrap(err, "RegisterDemoServiceHandlerFromEndpoint error")
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Viper.GetInt("server.http.port")), mux)
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

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Viper.GetInt("server.grpc.port")))
		if err != nil {
			panic(err)
		}

		s := grpc.NewServer()
		demo.RegisterDemoServiceServer(s, &server.Server{})

		if err = s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	if err := run(); err != nil {
		panic(err)
	}
}
