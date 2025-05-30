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

// ListClustersCmd holds the cmd flags
type ListClustersCmd struct {
	*flags.GlobalFlags
	Log log.Logger

	Host    string
	Project string
}

// NewListClustersCmd creates a new command
func NewListClustersCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &ListClustersCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Use:    "list-clusters",
		Short:  "List clusters",
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

func (cmd *ListClustersCmd) Run(ctx context.Context, devSpaceConfig *config.Config, provider *provider.ProviderConfig) error {
	opts := devSpaceConfig.ProviderOptions(provider.Name)
	opts[platform.ProjectEnv] = config.OptionValue{Value: cmd.Project}

	// ignore --debug because we tunnel json through stdio
	cmd.Log.SetLevel(logrus.InfoLevel)

	var buf bytes.Buffer
	err := clientimplementation.RunCommandWithBinaries(
		ctx,
		"listClusters",
		provider.Exec.Proxy.List.Clusters,
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
		return fmt.Errorf("list clusters with provider \"%s\": %w", provider.Name, err)
	}

	fmt.Println(buf.String())

	return nil
}
