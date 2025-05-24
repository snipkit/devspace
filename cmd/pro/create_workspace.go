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

// CreateWorkspaceCmd holds the cmd flags
type CreateWorkspaceCmd struct {
	*flags.GlobalFlags
	Log log.Logger

	Host     string
	Instance string
}

// NewCreateWorkspaceCmd creates a new command
func NewCreateWorkspaceCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &CreateWorkspaceCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Use:    "create-workspace",
		Short:  "Create workspace instance",
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
	c.Flags().StringVar(&cmd.Instance, "instance", "", "The workspace instance to create")
	_ = c.MarkFlagRequired("instance")

	return c
}

func (cmd *CreateWorkspaceCmd) Run(ctx context.Context, devSpaceConfig *config.Config, provider *provider.ProviderConfig) error {
	opts := devSpaceConfig.ProviderOptions(provider.Name)
	opts[platform.WorkspaceInstanceEnv] = config.OptionValue{Value: cmd.Instance}

	var buf bytes.Buffer
	// ignore --debug because we tunnel json through stdio
	cmd.Log.SetLevel(logrus.InfoLevel)

	err := clientimplementation.RunCommandWithBinaries(
		ctx,
		"createWorkspace",
		provider.Exec.Proxy.Create.Workspace,
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
		return fmt.Errorf("create workspace: %w", err)
	}

	fmt.Println(buf.String())

	return nil
}
