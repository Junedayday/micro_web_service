package server

import (
	"context"

	"github.com/Junedayday/micro_web_service/gen/idl/demo"
)

func (s *Server) Echo(ctx context.Context, req *demo.DemoRequest) (*demo.DemoResponse, error) {
	return &demo.DemoResponse{
		Code: 0,
	}, nil
}
