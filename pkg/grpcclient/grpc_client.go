package grpcclient

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/henderiw/mtls/apis/greeter/greeterpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	ctrl "sigs.k8s.io/controller-runtime"
)

type GRPCclient interface {
	Start(ctx context.Context) error
	Stop()
}

type grpcClient struct {
	config Config

	// logger
	l logr.Logger
	// cancel
	cancel context.CancelFunc
}

type Option func(GRPCclient)

func New(cfg Config, opts ...Option) GRPCclient {
	l := ctrl.Log.WithName("grpc server")
	cfg.setDefaults()
	r := &grpcClient{
		config: cfg,
		l:      l,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

func (r *grpcClient) Stop() {
	if r.cancel != nil {
		r.cancel()
		r.cancel = nil
	}
}

func (r *grpcClient) Start(ctx context.Context) error {
	r.l.Info("grpc client start",
		"address", r.config.Address,
		"certDir", r.config.CertDir,
		"certName", r.config.CertName,
		"keyName", r.config.KeyName,
		"caName", r.config.CaName,
	)
	tlsConfig, err := r.createTLSConfig(ctx)
	if err != nil {
		r.l.Error(err, "cannot create tls config")
		return err
	}
	conn, err := grpc.Dial(r.config.Address, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		r.l.Error(err, "cannot create grpc connection")
		return err
	}
	defer conn.Close()

	client := greeterpb.NewGreeterServiceClient(conn)

	for {
		resp, err := client.Hello(ctx, &greeterpb.HelloRequest{
			Name: "wim",
		})
		if err != nil {
			r.l.Error(err, "cannot get greeter response")
			return err
		}
		r.l.Info("hello", "resp", resp.Msg)

		time.Sleep(5 * time.Second)
	}
}
