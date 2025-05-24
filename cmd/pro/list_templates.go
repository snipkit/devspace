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

// ListTemplatesCmd holds the cmd flags
type ListTemplatesCmd struct {
	*flags.GlobalFlags
	Log log.Logger

	Host    string
	Project string
}

// NewListTemplatesCmd creates a new command
func NewListTemplatesCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &ListTemplatesCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Use:    "list-templates",
		Short:  "List templates",
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
	c.Flags().StringVar(&cmd.Project, "project", "", "The project to use")
	_ = c.MarkFlagRequired("project")

	return c
}

func (cmd *ListTemplatesCmd) Run(ctx context.Context, devSpaceConfig *config.Config, provider *provider.ProviderConfig) error {
	opts := devSpaceConfig.ProviderOptions(provider.Name)
	opts[platform.ProjectEnv] = config.OptionValue{Value: cmd.Project}

	// ignore --debug because we tunnel json through stdio
	cmd.Log.SetLevel(logrus.InfoLevel)
	var buf bytes.Buffer
	err := clientimplementation.RunCommandWithBinaries(
		ctx,
		"listTemplates",
		provider.Exec.Proxy.List.Templates,
		devSpaceConfig.DefaultContext,
		nil,
		nil,
		opts,
		provider,
		nil,
		nil,
		&buf,
		nil,
		cmd.Log)
	if err != nil {
		return fmt.Errorf("list templates with provider \"%s\": %w", provider.Name, err)
	}

	fmt.Println(buf.String())

	return nil
}
