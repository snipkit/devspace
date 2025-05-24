package json

import (
	"dev.khulnasoft.com/cmd/flags"
	"github.com/spf13/cobra"
)

// NewJSONCmd returns a new command
func NewJSONCmd(flags *flags.GlobalFlags) *cobra.Command {
	jsonCmd := &cobra.Command{
		Use:    "json",
		Short:  "DevSpace JSON Utility Commands",
		Hidden: true,
	}

	jsonCmd.AddCommand(NewGetCmd())
	return jsonCmd
}
