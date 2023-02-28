package main

import (
	"fmt"
	"os"

	"github.com/henderiw/mtls/pkg/run"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
	"k8s.io/component-base/cli"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	os.Exit(runMain())
}

// runMain does the initial setup in order to run funcbuilder. The return value from
// this function will be the exit code when funcbuilder terminates.
func runMain() int {
	var err error
	opts := zap.Options{
		Development: true,
		TimeEncoder: zapcore.ISO8601TimeEncoder,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	ctx := ctrl.SetupSignalHandler()

	// Enable commandline flags for klog.
	// logging will help in collecting debugging information from users
	klog.InitFlags(nil)

	cmd := run.GetMain(ctx)

	err = cli.RunNoErrOutput(cmd)
	if err != nil {
		return handleErr(cmd, err)
	}
	return 0
}

// handleErr takes care of printing an error message for a given error.
func handleErr(cmd *cobra.Command, err error) int {
	fmt.Fprintf(cmd.ErrOrStderr(), "%s \n", err.Error())
	return 1
}
