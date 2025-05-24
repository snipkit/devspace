package daemon

import (
	"context"
	"encoding/json"
	"fmt"

	platformdaemon "dev.khulnasoft.com/pkg/daemon/platform"

	"dev.khulnasoft.com/cmd/agent"
	"dev.khulnasoft.com/cmd/pro/completion"
	proflags "dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/config"
	providerpkg "dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
)

// StatusCmd holds the DevSpace daemon flags
type StatusCmd struct {
	*proflags.GlobalFlags

	Host string
	Log  log.Logger
}

// NewStatusCmd creates a new command
func NewStatusCmd(flags *proflags.GlobalFlags) *cobra.Command {
	cmd := &StatusCmd{
		GlobalFlags: flags,
		Log:         log.Default,
	}
	c := &cobra.Command{
		Use:   "status",
		Short: "Get the status of the daemon",
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
	_ = c.RegisterFlagCompletionFunc("host", func(rootCmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completion.GetPlatformHostSuggestions(rootCmd, cmd.Context, cmd.Provider, args, toComplete, cmd.Owner, cmd.Log)
	})

	return c
}

func (cmd *StatusCmd) Run(ctx context.Context, devSpaceConfig *config.Config, provider *providerpkg.ProviderConfig) error {
	status, err := platformdaemon.NewLocalClient(provider.Name).Status(ctx, cmd.Debug)
	if err != nil {
		return err
	}
	out, err := json.Marshal(status)
	if err != nil {
		return err
	}

	fmt.Print(string(out))

	return nil
}
