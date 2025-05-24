package pro

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	proflags "dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"dev.khulnasoft.com/log/table"
	"github.com/spf13/cobra"
)

// ListCmd holds the list cmd flags
type ListCmd struct {
	proflags.GlobalFlags

	Output string
	Login  bool
}

// NewListCmd creates a new command
func NewListCmd(flags *proflags.GlobalFlags) *cobra.Command {
	cmd := &ListCmd{
		GlobalFlags: *flags,
	}
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available DevSpace Pro instances",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.Run(context.Background())
		},
	}

	listCmd.Flags().StringVar(&cmd.Output, "output", "plain", "The output format to use. Can be json or plain")
	listCmd.Flags().BoolVar(&cmd.Login, "login", false, "Check if the user is logged into the pro instance")
	return listCmd
}

// Run runs the command logic
func (cmd *ListCmd) Run(ctx context.Context) error {
	devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
	if err != nil {
		return err
	}

	proInstances, err := workspace.ListProInstances(devSpaceConfig, log.Default)
	if err != nil {
		return err
	}

	if cmd.Output == "plain" {
		tableEntries := [][]string{}
		for _, proInstance := range proInstances {
			entry := []string{
				proInstance.Host,
				proInstance.Provider,
				time.Since(proInstance.CreationTimestamp.Time).Round(1 * time.Second).String(),
			}
			if cmd.Login {
				err = checkLogin(ctx, devSpaceConfig, proInstance)
				entry = append(entry, fmt.Sprintf("%t", err == nil))
			}

			tableEntries = append(tableEntries, entry)
		}
		sort.SliceStable(tableEntries, func(i, j int) bool {
			return tableEntries[i][0] < tableEntries[j][0]
		})

		tableHeaders := []string{
			"Host",
			"Provider",
			"Age",
		}
		if cmd.Login {
			tableHeaders = append(tableHeaders, "Authenticated")
		}

		table.PrintTable(log.Default, tableHeaders, tableEntries)
	} else if cmd.Output == "json" {
		tableEntries := []*proTableEntry{}
		for _, proInstance := range proInstances {
			entry := &proTableEntry{
				ProInstance:  proInstance,
				Context:      devSpaceConfig.DefaultContext,
				Capabilities: getCapabilities(ctx, devSpaceConfig, proInstance, log.Discard),
			}
			if cmd.Login {
				err = checkLogin(ctx, devSpaceConfig, proInstance)
				isAuthenticated := err == nil
				entry.Authenticated = &isAuthenticated
			}

			tableEntries = append(tableEntries, entry)
		}

		sort.SliceStable(tableEntries, func(i, j int) bool {
			return tableEntries[i].Host < tableEntries[j].Host
		})
		out, err := json.Marshal(tableEntries)
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	} else {
		return fmt.Errorf("unexpected output format, choose either json or plain. Got %s", cmd.Output)
	}

	return nil
}

type proTableEntry struct {
	*provider.ProInstance

	Authenticated *bool        `json:"authenticated,omitempty"`
	Context       string       `json:"context,omitempty"`
	Capabilities  []Capability `json:"capabilities,omitempty"`
}

type Capability string

var (
	capabilityDaemon         Capability = "daemon"
	capabilityHealthCheck    Capability = "health-check"
	capabilityUpdateProvider Capability = "update-provider"
)

func checkLogin(ctx context.Context, devSpaceConfig *config.Config, proInstance *provider.ProInstance) error {
	// for every pro instance, check auth status by calling login
	if err := login(ctx, devSpaceConfig, proInstance.Host, proInstance.Provider, "", true, false, log.Default); err != nil {
		return fmt.Errorf("not logged into %s", proInstance.Host)
	}

	return nil
}

func getCapabilities(ctx context.Context, devSpaceConfig *config.Config, proInstance *provider.ProInstance, log log.Logger) []Capability {
	capabilities := []Capability{}
	provider, err := workspace.FindProvider(devSpaceConfig, proInstance.Provider, log)
	if err != nil {
		return capabilities
	}

	if provider.Config.HasHealthCheck() {
		capabilities = append(capabilities, capabilityHealthCheck)
		capabilities = append(capabilities, capabilityUpdateProvider)
	}

	if provider.Config.IsDaemonProvider() {
		capabilities = append(capabilities, capabilityDaemon)
	}

	return capabilities
}
