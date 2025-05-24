package machine

import (
	"context"
	"fmt"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/client"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
)

// DeleteCmd holds the configuration
type DeleteCmd struct {
	*flags.GlobalFlags

	GracePeriod string
	Force       bool
}

// NewDeleteCmd creates a new destroy command
func NewDeleteCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &DeleteCmd{
		GlobalFlags: flags,
	}
	deleteCmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Deletes an existing machine",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.Run(context.Background(), args)
		},
	}

	deleteCmd.Flags().StringVar(&cmd.GracePeriod, "grace-period", "", "The amount of time to give the command to delete the workspace")
	deleteCmd.Flags().BoolVar(&cmd.Force, "force", false, "Delete workspace even if it is not found remotely anymore")
	return deleteCmd
}

// Run runs the command logic
func (cmd *DeleteCmd) Run(ctx context.Context, args []string) error {
	devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
	if err != nil {
		return err
	}

	machineClient, err := workspace.GetMachine(devSpaceConfig, args, log.Default)
	if err != nil {
		return err
	}

	// check if there are workspaces that still use this machine
	workspaces, err := workspace.List(ctx, devSpaceConfig, false, platform.SelfOwnerFilter, log.Default)
	if err != nil {
		return err
	}

	// search for workspace that uses this machine
	for _, workspace := range workspaces {
		if workspace.Machine.ID == machineClient.Machine() {
			return fmt.Errorf("cannot delete machine '%s', because workspace '%s' is still using it. Please delete the workspace '%s' before deleting the machine", workspace.Machine.ID, workspace.ID, workspace.ID)
		}
	}

	err = machineClient.Delete(ctx, client.DeleteOptions{
		Force:       cmd.Force,
		GracePeriod: cmd.GracePeriod,
	})
	if err != nil {
		return err
	}

	return nil
}
