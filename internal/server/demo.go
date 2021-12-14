package server

import (
	"context"

	"github.com/Junedayday/micro_web_service/gen/idl/demo"
)

func (s *Server) Demo(ctx context.Context, req *demo.DemoRequest) (*demo.DemoResponse, error) {
	return &demo.DemoResponse{}, nil
}

func (s *Server) Empty(ctx context.Context, req *demo.EmptyRequest) (*demo.EmptyResponse, error) {
	return &demo.EmptyResponse{}, nil
}
