package pro

import (
	"context"
	"fmt"
	"os"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/cmd/pro/add"
	"dev.khulnasoft.com/cmd/pro/daemon"
	proflags "dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/cmd/pro/provider"
	"dev.khulnasoft.com/cmd/pro/reset"
	"dev.khulnasoft.com/pkg/config"
	providerpkg "dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewProCmd returns a new command
func NewProCmd(flags *flags.GlobalFlags, streamLogger *log.StreamLogger) *cobra.Command {
	globalFlags := &proflags.GlobalFlags{GlobalFlags: flags}
	proCmd := &cobra.Command{
		Use:           "pro",
		Short:         "DevSpace Pro commands",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		PersistentPreRunE: func(c *cobra.Command, _ []string) error {
			globalFlags = proflags.SetGlobalFlags(c.PersistentFlags())
			if flags.Silent {
				streamLogger.SetLevel(logrus.FatalLevel)
			}
			if flags.Debug {
				streamLogger.SetLevel(logrus.DebugLevel)
			}

			if os.Getenv("DEVSPACE_DEBUG") == "true" {
				log.Default.SetLevel(logrus.DebugLevel)
			}
			if flags.LogOutput == "json" {
				log.Default.SetFormat(log.JSONFormat)
			}

			return nil
		},
	}

	proCmd.AddCommand(NewLoginCmd(globalFlags))
	proCmd.AddCommand(NewListCmd(globalFlags))
	proCmd.AddCommand(NewDeleteCmd(globalFlags))
	proCmd.AddCommand(NewImportCmd(globalFlags))
	proCmd.AddCommand(NewStartCmd(globalFlags))
	proCmd.AddCommand(NewRebuildCmd(globalFlags))
	proCmd.AddCommand(NewSleepCmd(globalFlags))
	proCmd.AddCommand(NewWakeupCmd(globalFlags))
	proCmd.AddCommand(reset.NewResetCmd(globalFlags))
	proCmd.AddCommand(provider.NewProProviderCmd(globalFlags))
	proCmd.AddCommand(daemon.NewCmd(globalFlags))
	proCmd.AddCommand(add.NewAddCmd(globalFlags))
	proCmd.AddCommand(NewWatchWorkspacesCmd(globalFlags))
	proCmd.AddCommand(NewSelfCmd(globalFlags))
	proCmd.AddCommand(NewVersionCmd(globalFlags))
	proCmd.AddCommand(NewListProjectsCmd(globalFlags))
	proCmd.AddCommand(NewListWorkspacesCmd(globalFlags))
	proCmd.AddCommand(NewListTemplatesCmd(globalFlags))
	proCmd.AddCommand(NewListClustersCmd(globalFlags))
	proCmd.AddCommand(NewCreateWorkspaceCmd(globalFlags))
	proCmd.AddCommand(NewUpdateWorkspaceCmd(globalFlags))
	proCmd.AddCommand(NewCheckHealthCmd(globalFlags))
	proCmd.AddCommand(NewCheckUpdateCmd(globalFlags))
	proCmd.AddCommand(NewUpdateProviderCmd(globalFlags))
	return proCmd
}

func findProProvider(ctx context.Context, context, provider, host string, log log.Logger) (*config.Config, *providerpkg.ProviderConfig, error) {
	devSpaceConfig, err := config.LoadConfig(context, provider)
	if err != nil {
		return nil, nil, err
	}

	pCfg, err := workspace.ProviderFromHost(ctx, devSpaceConfig, host, log)
	if err != nil {
		return devSpaceConfig, nil, fmt.Errorf("load provider: %w", err)
	}

	return devSpaceConfig, pCfg, nil
}
