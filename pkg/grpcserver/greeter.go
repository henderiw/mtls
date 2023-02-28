package grpcserver

import (
	"context"

	"github.com/henderiw/mtls/apis/greeter/greeterpb"
)

func (r *grpcServer) Hello(ctx context.Context, req *greeterpb.HelloRequest) (*greeterpb.HelloResponse, error) {
	r.l.Info("greeter hello", "req name", req.GetName())
	respdata := "Hello," + req.GetName()
	return &greeterpb.HelloResponse{Msg: respdata}, nil
}
