package pro

import (
	"bytes"
	"context"
	"fmt"

	"dev.khulnasoft.com/cmd/agent"
	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/client/clientimplementation"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CheckHealthCmd holds the cmd flags
type CheckHealthCmd struct {
	*flags.GlobalFlags
	Log log.Logger

	Host string
}

// NewCheckHealthCmd creates a new command
func NewCheckHealthCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &CheckHealthCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Use:    "check-health",
		Short:  "Check platform health",
		Hidden: true,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			devSpaceConfig, provider, err := findProProvider(cobraCmd.Context(), cmd.Context, cmd.Provider, cmd.Host, cmd.Log)
			if err != nil {
				return err
			}

			return cmd.Run(cobraCmd.Context(), devSpaceConfig, provider)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			root := cmd.Root()
			if root == nil {
				return
			}
			if root.Annotations == nil {
				root.Annotations = map[string]string{}
			}
			// Don't print debug message
			root.Annotations[agent.AgentExecutedAnnotation] = "true"
		},
	}

	c.Flags().StringVar(&cmd.Host, "host", "", "The pro instance to use")
	_ = c.MarkFlagRequired("host")

	return c
}

func (cmd *CheckHealthCmd) Run(ctx context.Context, devSpaceConfig *config.Config, provider *provider.ProviderConfig) error {
	var buf bytes.Buffer
	// ignore --debug because we tunnel json through stdio
	cmd.Log.SetLevel(logrus.InfoLevel)

	err := clientimplementation.RunCommandWithBinaries(
		ctx,
		"health",
		provider.Exec.Proxy.Health,
		devSpaceConfig.DefaultContext,
		nil,
		nil,
		devSpaceConfig.ProviderOptions(provider.Name),
		provider,
		nil,
		nil,
		&buf,
		cmd.Log.Writer(logrus.ErrorLevel, true),
		cmd.Log)
	if err != nil {
		return fmt.Errorf("check health with provider \"%s\": %w", provider.Name, err)
	}

	fmt.Println(buf.String())

	return nil
}
