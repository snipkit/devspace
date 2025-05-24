package provider

import (
	"context"
	"fmt"
	"io"
	"os"

	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteCmd holds the cmd flags
type DeleteCmd struct {
	*flags.GlobalFlags

	Log log.Logger
}

// NewDeleteCmd creates a new command
func NewDeleteCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &DeleteCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Hidden: true,
		Use:    "delete",
		Short:  "Runs delete on a workspace",
		Args:   cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run(cobraCmd.Context(), os.Stdin, os.Stdout, os.Stderr)
		},
	}

	return c
}

func (cmd *DeleteCmd) Run(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	baseClient, err := client.InitClientFromPath(ctx, cmd.Config)
	if err != nil {
		return err
	}

	info, err := platform.GetWorkspaceInfoFromEnv()
	if err != nil {
		return err
	}
	workspace, err := platform.FindInstanceInProject(ctx, baseClient, info.UID, info.ProjectName)
	if err != nil {
		return err
	} else if workspace == nil {
		return fmt.Errorf("couldn't find workspace")
	}

	managementClient, err := baseClient.Management()
	if err != nil {
		return err
	}

	err = managementClient.Loft().ManagementV1().DevSpaceWorkspaceInstances(workspace.Namespace).Delete(ctx, workspace.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("delete workspace: %w", err)
	}

	return nil
}
