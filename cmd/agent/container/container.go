package container

import (
	"dev.khulnasoft.com/cmd/flags"
	"github.com/spf13/cobra"
)

// NewContainerCmd returns a new command
func NewContainerCmd(flags *flags.GlobalFlags) *cobra.Command {
	containerCmd := &cobra.Command{
		Use:   "container",
		Short: "Container commands",
	}

	containerCmd.AddCommand(NewSetupContainerCmd(flags))
	containerCmd.AddCommand(NewDaemonCmd())
	containerCmd.AddCommand(NewVSCodeAsyncCmd())
	containerCmd.AddCommand(NewOpenVSCodeAsyncCmd())
	containerCmd.AddCommand(NewCredentialsServerCmd(flags))
	containerCmd.AddCommand(NewSetupKhulnasoftPlatformAccessCmd(flags))
	containerCmd.AddCommand(NewSSHServerCmd(flags))
	return containerCmd
}
