package get

import (
	"dev.khulnasoft.com/cmd/pro/flags"
	"github.com/spf13/cobra"
)

// NewCmd creates a new cobra command
func NewCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	c := &cobra.Command{
		Use:    "get",
		Short:  "DevSpace Pro Provider get commands",
		Args:   cobra.NoArgs,
		Hidden: true,
	}

	c.AddCommand(NewWorkspaceCmd(globalFlags))
	c.AddCommand(NewSelfCmd(globalFlags))
	c.AddCommand(NewVersionCmd(globalFlags))

	return c
}
