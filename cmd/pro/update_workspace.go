package pro

import (
	"bytes"
	"context"
	"fmt"

	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/client/clientimplementation"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// UpdateWorkspaceCmd holds the cmd flags
type UpdateWorkspaceCmd struct {
	*flags.GlobalFlags
	Log log.Logger

	Host     string
	Instance string
}

// NewUpdateWorkspaceCmd creates a new command
func NewUpdateWorkspaceCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &UpdateWorkspaceCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Use:    "update-workspace",
		Short:  "Update workspace instance",
		Hidden: true,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			devSpaceConfig, provider, err := findProProvider(cobraCmd.Context(), cmd.Context, cmd.Provider, cmd.Host, cmd.Log)
			if err != nil {
				return err
			}

			return cmd.Run(cobraCmd.Context(), devSpaceConfig, provider)
		},
	}

	c.Flags().StringVar(&cmd.Host, "host", "", "The pro instance to use")
	_ = c.MarkFlagRequired("host")
	c.Flags().StringVar(&cmd.Instance, "instance", "", "The workspace instance to update")
	_ = c.MarkFlagRequired("instance")

	return c
}

func (cmd *UpdateWorkspaceCmd) Run(ctx context.Context, devSpaceConfig *config.Config, provider *provider.ProviderConfig) error {
	opts := devSpaceConfig.ProviderOptions(provider.Name)
	opts[platform.WorkspaceInstanceEnv] = config.OptionValue{Value: cmd.Instance}

	var buf bytes.Buffer
	// ignore --debug because we tunnel json through stdio
	cmd.Log.SetLevel(logrus.InfoLevel)

	err := clientimplementation.RunCommandWithBinaries(
		ctx,
		"updateWorkspace",
		provider.Exec.Proxy.Update.Workspace,
		devSpaceConfig.DefaultContext,
		nil,
		nil,
		opts,
		provider,
		nil,
		nil,
		&buf,
		cmd.Log.ErrorStreamOnly().Writer(logrus.ErrorLevel, true),
		cmd.Log)
	if err != nil {
		return fmt.Errorf("update workspace with provider \"%s\": %w", provider.Name, err)
	}

	fmt.Println(buf.String())

	return nil
}
