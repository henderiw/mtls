package grpcclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
)

func (r *grpcClient) createTLSConfig(ctx context.Context) (*tls.Config, error) {
	caFile := filepath.Join(r.config.CertDir, r.config.CaName)
	certFile := filepath.Join(r.config.CertDir, r.config.CertName)
	keyFile := filepath.Join(r.config.CertDir, r.config.KeyName)

	/*
		creds, err := credentials.NewClientTLSFromFile(certPath, "")
		if err != nil {
			return nil, err
		}
	*/

	ca, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client CA cert: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("cannot add ca cert to the ca pool")
	}

	// client does not require certwatcher since the client
	// needs to reconnect and need to read the latest CERT
	creds, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{creds},
		RootCAs: caPool, // regular tls validation - check if the server is valid (rootCA)
	}
	return tlsConfig, nil
}
