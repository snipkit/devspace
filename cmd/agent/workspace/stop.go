package workspace

import (
	"context"
	"fmt"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/agent"
	provider2 "dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// StopCmd holds the cmd flags
type StopCmd struct {
	*flags.GlobalFlags

	WorkspaceInfo string
}

// NewStopCmd creates a new command
func NewStopCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &StopCmd{
		GlobalFlags: flags,
	}
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops a workspace on the remote server",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmd.Run(context.Background())
		},
	}
	stopCmd.Flags().StringVar(&cmd.WorkspaceInfo, "workspace-info", "", "The workspace info")
	_ = stopCmd.MarkFlagRequired("workspace-info")
	return stopCmd
}

func (cmd *StopCmd) Run(ctx context.Context) error {
	// get workspace
	shouldExit, workspaceInfo, err := agent.WriteWorkspaceInfo(cmd.WorkspaceInfo, log.Default.ErrorStreamOnly())
	if err != nil {
		return fmt.Errorf("error parsing workspace info: %w", err)
	} else if shouldExit {
		return nil
	}

	// stop docker container
	err = stopContainer(ctx, workspaceInfo, log.Default)
	if err != nil {
		return errors.Wrap(err, "stop container")
	}

	return nil
}

func stopContainer(ctx context.Context, workspaceInfo *provider2.AgentWorkspaceInfo, log log.Logger) error {
	log.Debugf("Stopping DevSpace container...")
	runner, err := CreateRunner(workspaceInfo, log)
	if err != nil {
		return err
	}

	err = runner.Stop(ctx)
	if err != nil {
		return err
	}
	log.Debugf("Successfully stopped DevSpace container")

	return nil
}
