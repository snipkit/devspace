package list

import (
	"dev.khulnasoft.com/cmd/pro/flags"
	"github.com/spf13/cobra"
)

// NewCmd creates a new cobra command
func NewCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	c := &cobra.Command{
		Use:    "list",
		Short:  "DevSpace Pro Provider list commands",
		Args:   cobra.NoArgs,
		Hidden: true,
	}

	c.AddCommand(NewWorkspacesCmd(globalFlags))
	c.AddCommand(NewProjectsCmd(globalFlags))
	c.AddCommand(NewTemplatesCmd(globalFlags))
	c.AddCommand(NewClustersCmd(globalFlags))

	return c
}
