package client

import (
	"context"

	"github.com/henderiw/mtls/pkg/grpcclient"
	"github.com/spf13/cobra"
)

// NewRunner returns a command runner.
func NewRunner(ctx context.Context, parent string) *Runner {
	r := &Runner{
		ctx:     ctx,
		errChan: make(chan error),
	}
	c := &cobra.Command{
		Use:      "client",
		PreRunE:  r.prerunE,
		RunE:     r.runE,
		PostRunE: r.postrunE,
	}

	c.Flags().BoolVar(&r.local, "local", false, "run using local certificates")

	r.Command = c
	return r
}

func NewCommand(ctx context.Context, parent string) *cobra.Command {
	return NewRunner(ctx, parent).Command
}

type Runner struct {
	Command *cobra.Command
	ctx     context.Context
	errChan chan error
	local   bool
}

func (r *Runner) prerunE(c *cobra.Command, args []string) error {

	return nil
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	cfg := grpcclient.Config{
		Address: "localhost:8888",
	}
	if r.local {
		cfg.CertDir = "certs"
		cfg.CertName = "client.crt"
		cfg.KeyName = "client.key"
		cfg.CaName = "ca.crt"
	}

	client := grpcclient.New(cfg)
	go func() {
		if err := client.Start(r.ctx); err != nil {
			r.errChan <- err
		}
	}()

	// block until cancelled or err occurs
	for {
		select {
		case <-r.ctx.Done():
			// We are done
			return nil
		case err := <-r.errChan:
			// Error starting or during start
			return err
		}
	}
}

func (r *Runner) postrunE(c *cobra.Command, args []string) error {
	return nil
}
