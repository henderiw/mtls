package grpcclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/controller-runtime/pkg/certwatcher"
)

func (r *grpcClient) createTLSConfig(ctx context.Context) (*tls.Config, error) {
	caPath := filepath.Join(r.config.CertDir, r.config.CaName)
	certPath := filepath.Join(r.config.CertDir, r.config.CertName)
	keyPath := filepath.Join(r.config.CertDir, r.config.KeyName)

	ca, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read client CA cert: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("cannot add ca cert to the ca pool")
	}

	certWatcher, err := certwatcher.New(certPath, keyPath)
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
		//RootCAs: caPool, // regular tls
	}
	return tlsConfig, nil
}
