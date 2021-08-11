package main

import (
	"context"
	"flag"
	"net"
	"net/http"

	"github.com/Junedayday/micro_web_service/gen/idl/demo"
	"github.com/Junedayday/micro_web_service/internal/server"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", ":9090", "gRPC server endpoint")
)

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := demo.RegisterDemoServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return errors.Wrap(err, "RegisterDemoServiceHandlerFromEndpoint error")
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", mux)
}

func main() {
	flag.Parse()

	go func() {
		lis, err := net.Listen("tcp", ":9090")
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
