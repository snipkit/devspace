package machine

import (
	"context"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/client"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
)

// StartCmd holds the configuration
type StartCmd struct {
	*flags.GlobalFlags
}

// NewStartCmd creates a new destroy command
func NewStartCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &StartCmd{
		GlobalFlags: flags,
	}
	startCmd := &cobra.Command{
		Use:   "start [name]",
		Short: "Starts an existing machine",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.Run(context.Background(), args)
		},
	}

	return startCmd
}

// Run runs the command logic
func (cmd *StartCmd) Run(ctx context.Context, args []string) error {
	devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
	if err != nil {
		return err
	}

	machineClient, err := workspace.GetMachine(devSpaceConfig, args, log.Default)
	if err != nil {
		return err
	}

	err = machineClient.Start(ctx, client.StartOptions{})
	if err != nil {
		return err
	}

	return nil
}
