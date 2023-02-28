package grpcserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"sigs.k8s.io/controller-runtime/pkg/certwatcher"
)

func (r *grpcServer) serverOpts(ctx context.Context) ([]grpc.ServerOption, error) {
	if r.config.Insecure {
		return []grpc.ServerOption{
			grpc.Creds(insecure.NewCredentials()),
		}, nil
	}

	tlsConfig, err := r.createTLSConfig(ctx)
	if err != nil {
		return nil, err
	}
	return []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		//grpc.UnaryInterceptor(r.middleFunc), //mtls
	}, nil

}

func (r *grpcServer) createTLSConfig(ctx context.Context) (*tls.Config, error) {
	caFile := filepath.Join(r.config.CertDir, r.config.CaName)
	certFile := filepath.Join(r.config.CertDir, r.config.CertName)
	keyFile := filepath.Join(r.config.CertDir, r.config.KeyName)

	ca, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client CA cert: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("cannot add ca cert to the ca pool")
	}

	certWatcher, err := certwatcher.New(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := certWatcher.Start(ctx); err != nil {
			r.l.Info("certificate watcher", "error", err)
		}
	}()

	tlsConfig := &tls.Config{
		GetCertificate: certWatcher.GetCertificate,
		ClientAuth:     tls.RequireAndVerifyClientCert, // mtls
		ClientCAs:      caPool,                         // mtls
		RootCAs:        caPool,                         // regular tls
	}
	return tlsConfig, nil
}

/*
func (r *grpcServer) middleFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// get client tls info
	if p, ok := peer.FromContext(ctx); ok {
		if mtls, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			for _, item := range mtls.State.PeerCertificates {
				r.l.Info("request certificate", "subject", item.Subject)
			}
		}
	}
	return handler(ctx, req)
}
*/
