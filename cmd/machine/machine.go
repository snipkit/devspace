package machine

import (
	"dev.khulnasoft.com/cmd/flags"
	"github.com/spf13/cobra"
)

// NewMachineCmd returns a new root command
func NewMachineCmd(flags *flags.GlobalFlags) *cobra.Command {
	machineCmd := &cobra.Command{
		Use:   "machine",
		Short: "DevSpace Machine commands",
	}

	machineCmd.AddCommand(NewListCmd(flags))
	machineCmd.AddCommand(NewSSHCmd(flags))
	machineCmd.AddCommand(NewStopCmd(flags))
	machineCmd.AddCommand(NewStartCmd(flags))
	machineCmd.AddCommand(NewStatusCmd(flags))
	machineCmd.AddCommand(NewDeleteCmd(flags))
	machineCmd.AddCommand(NewCreateCmd(flags))
	machineCmd.AddCommand(NewInspectCmd(flags))
	return machineCmd
}
