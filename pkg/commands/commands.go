package commands

import (
	"context"

	"github.com/henderiw/mtls/pkg/commands/client"
	"github.com/henderiw/mtls/pkg/commands/server"
	"github.com/spf13/cobra"
)

// GetMTLSCommands returns the set of mtls commands to be registered
func GetMTLSCommands(ctx context.Context, name string) []*cobra.Command {
	var c []*cobra.Command
	clientCmd := client.NewCommand(ctx, name)
	serverCmd := server.NewCommand(ctx, name)

	c = append(c, clientCmd, serverCmd)
	return c
}
