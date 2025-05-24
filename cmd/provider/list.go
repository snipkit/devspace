package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/types"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"dev.khulnasoft.com/log/table"
	"github.com/spf13/cobra"
)

// ListCmd holds the list cmd flags
type ListCmd struct {
	flags.GlobalFlags

	Output string
}

// NewListCmd creates a new command
func NewListCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &ListCmd{
		GlobalFlags: *flags,
	}
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available providers",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.Run(context.Background())
		},
	}

	listCmd.Flags().StringVar(&cmd.Output, "output", "plain", "The output format to use. Can be json or plain")
	return listCmd
}

type ProviderWithDefault struct {
	workspace.ProviderWithOptions `json:",inline"`

	Default bool `json:"default,omitempty"`
}

// Run runs the command logic
func (cmd *ListCmd) Run(ctx context.Context) error {
	devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
	if err != nil {
		return err
	}

	providers, err := workspace.LoadAllProviders(devSpaceConfig, log.Default.ErrorStreamOnly())
	if err != nil {
		return err
	}

	configuredProviders := devSpaceConfig.Current().Providers
	if configuredProviders == nil {
		configuredProviders = map[string]*config.ProviderConfig{}
	}

	if cmd.Output == "plain" {
		tableEntries := [][]string{}
		for _, entry := range providers {
			tableEntries = append(tableEntries, []string{
				entry.Config.Name,
				entry.Config.Version,
				strconv.FormatBool(devSpaceConfig.Current().DefaultProvider == entry.Config.Name),
				strconv.FormatBool(entry.State != nil && entry.State.Initialized),
				entry.Config.Description,
			})
		}
		sort.SliceStable(tableEntries, func(i, j int) bool {
			return tableEntries[i][0] < tableEntries[j][0]
		})

		table.PrintTable(log.Default, []string{
			"Name",
			"Version",
			"Default",
			"Initialized",
			"Description",
		}, tableEntries)
	} else if cmd.Output == "json" {
		retMap := map[string]ProviderWithDefault{}
		for k, entry := range providers {
			var dynamicOptions map[string]*types.Option
			if configuredProviders[entry.Config.Name] != nil {
				dynamicOptions = configuredProviders[entry.Config.Name].DynamicOptions
			}

			srcOptions := MergeDynamicOptions(entry.Config.Options, dynamicOptions)
			entry.Config.Options = srcOptions
			retMap[k] = ProviderWithDefault{
				ProviderWithOptions: *entry,
				Default:             devSpaceConfig.Current().DefaultProvider == entry.Config.Name,
			}
		}

		out, err := json.MarshalIndent(retMap, "", "  ")
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	} else {
		return fmt.Errorf("unexpected output format, choose either json or plain. Got %s", cmd.Output)
	}

	return nil
}
