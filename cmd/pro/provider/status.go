package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	storagev1 "dev.khulnasoft.com/api/pkg/apis/storage/v1"
	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/pkg/platform/remotecommand"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
)

// StatusCmd holds the cmd flags
type StatusCmd struct {
	*flags.GlobalFlags

	Log log.Logger
}

// NewStatusCmd creates a new command
func NewStatusCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &StatusCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Hidden: true,
		Use:    "status",
		Short:  "Runs status on a workspace",
		Args:   cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run(cobraCmd.Context(), os.Stdin, os.Stdout, os.Stderr)
		},
	}

	return c
}

func (cmd *StatusCmd) Run(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
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
		out, err := json.Marshal(&storagev1.WorkspaceStatusResult{
			ID:       os.Getenv(platform.WorkspaceIDEnv),
			Context:  os.Getenv(platform.WorkspaceContextEnv),
			State:    string(storagev1.WorkspaceStatusNotFound),
			Provider: os.Getenv(platform.WorkspaceProviderEnv),
		})
		if err != nil {
			return err
		}

		fmt.Println(string(out))
		return nil
	}

	conn, err := platform.DialInstance(baseClient, workspace, "getstatus", platform.OptionsFromEnv("DEVSPACE_FLAGS_STATUS"), cmd.Log)
	if err != nil {
		return err
	}

	_, err = remotecommand.ExecuteConn(ctx, conn, stdin, stdout, stderr, cmd.Log.ErrorStreamOnly())
	if err != nil {
		return fmt.Errorf("error executing: %w", err)
	}

	return nil
}
