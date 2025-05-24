package context

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
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
		Short:   "List DevSpace contexts",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.Run(context.Background())
		},
	}

	listCmd.Flags().StringVar(&cmd.Output, "output", "plain", "The output format to use. Can be json or plain")
	return listCmd
}

type ContextWithDefault struct {
	Name string `json:"name,omitempty"`

	Default bool `json:"default,omitempty"`
}

// Run runs the command logic
func (cmd *ListCmd) Run(ctx context.Context) error {
	devSpaceConfig, err := config.LoadConfig(cmd.Context, cmd.Provider)
	if err != nil {
		return err
	}

	if cmd.Output == "plain" {
		tableEntries := [][]string{}
		for contextName := range devSpaceConfig.Contexts {
			tableEntries = append(tableEntries, []string{
				contextName,
				strconv.FormatBool(devSpaceConfig.DefaultContext == contextName),
			})
		}
		sort.SliceStable(tableEntries, func(i, j int) bool {
			return tableEntries[i][0] < tableEntries[j][0]
		})

		table.PrintTable(log.Default, []string{
			"Name",
			"Default",
		}, tableEntries)
	} else if cmd.Output == "json" {
		ides := []ContextWithDefault{}
		for contextName := range devSpaceConfig.Contexts {
			ides = append(ides, ContextWithDefault{
				Name:    contextName,
				Default: devSpaceConfig.DefaultContext == contextName,
			})
		}

		out, err := json.MarshalIndent(ides, "", "  ")
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	} else {
		return fmt.Errorf("unexpected output format, choose either json or plain. Got %s", cmd.Output)
	}

	return nil
}
