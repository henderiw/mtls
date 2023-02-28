package grpcserver

import (
	"context"
	"net"
	"sync"

	"github.com/go-logr/logr"
	"github.com/henderiw/mtls/apis/greeter/greeterpb"
	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"
	ctrl "sigs.k8s.io/controller-runtime"
)

type GRPCserver interface {
	Start(ctx context.Context) error
	Stop()
}

type grpcServer struct {
	config Config
	greeterpb.UnimplementedGreeterServiceServer

	// logger
	l logr.Logger

	//Handlers

	// cached certificate
	m   *sync.Mutex
	sem *semaphore.Weighted

	cancel context.CancelFunc
}

type Option func(GRPCserver)

func New(c Config, opts ...Option) GRPCserver {
	l := ctrl.Log.WithName("grpc server")
	c.setDefaults()
	s := &grpcServer{
		config: c,
		sem:    semaphore.NewWeighted(c.MaxRPC),
		m:      new(sync.Mutex),
		l:      l,
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

func (r *grpcServer) Stop() {
	if r.cancel != nil {
		r.cancel()
		r.cancel = nil
	}
}

func (r *grpcServer) Start(ctx context.Context) error {
	r.l.Info("grpc server start",
		"address", r.config.Address,
		"certDir", r.config.CertDir,
		"certName", r.config.CertName,
		"keyName", r.config.KeyName,
		"caName", r.config.CaName,
	)
	ctx, cancel := context.WithCancel(ctx)
	r.cancel = cancel

	l, err := net.Listen("tcp", r.config.Address)
	if err != nil {
		return errors.Wrap(err, "cannot listen")
	}
	opts, err := r.serverOpts(ctx)
	if err != nil {
		return err
	}

	// create a gRPC server object
	s := grpc.NewServer(opts...)

	greeterpb.RegisterGreeterServiceServer(s, r)

	r.l.Info("starting grpc server...")
	if err := s.Serve(l); err != nil {
		r.l.Info("gRPC serve failed", "error", err)
		return err
	}
	return nil
}
