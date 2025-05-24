package helper

import (
	"dev.khulnasoft.com/cmd/agent"
	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/cmd/helper/http"
	"dev.khulnasoft.com/cmd/helper/json"
	"dev.khulnasoft.com/cmd/helper/strings"
	"github.com/spf13/cobra"
)

// NewHelperCmd returns a new command
func NewHelperCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	helperCmd := &cobra.Command{
		Use:   "helper",
		Short: "DevSpace Utility Commands",
		PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
			return agent.AgentPersistentPreRunE(cobraCmd, args, globalFlags)
		},
		Hidden: true,
	}

	helperCmd.AddCommand(http.NewHTTPCmd(globalFlags))
	helperCmd.AddCommand(json.NewJSONCmd(globalFlags))
	helperCmd.AddCommand(strings.NewStringsCmd(globalFlags))
	helperCmd.AddCommand(NewSSHServerCmd(globalFlags))
	helperCmd.AddCommand(NewGetWorkspaceNameCmd(globalFlags))
	helperCmd.AddCommand(NewGetWorkspaceUIDCmd(globalFlags))
	helperCmd.AddCommand(NewGetWorkspaceConfigCommand(globalFlags))
	helperCmd.AddCommand(NewGetProviderNameCmd(globalFlags))
	helperCmd.AddCommand(NewCheckProviderUpdateCmd(globalFlags))
	helperCmd.AddCommand(NewSSHClientCmd())
	helperCmd.AddCommand(NewShellCmd())
	helperCmd.AddCommand(NewSSHGitCloneCmd())
	helperCmd.AddCommand(NewFleetServerCmd(globalFlags))
	helperCmd.AddCommand(NewDockerCredentialsHelperCmd(globalFlags))
	helperCmd.AddCommand(NewGetImageCmd(globalFlags))
	return helperCmd
}
