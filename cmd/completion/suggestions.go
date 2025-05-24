package completion

import (
	"strings"

	"github.com/spf13/cobra"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
)

func RegisterFlagCompletionFuns(rootCmd *cobra.Command, globalFlags *flags.GlobalFlags) error {
	if err := rootCmd.RegisterFlagCompletionFunc("provider", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return GetProviderSuggestions(rootCmd, globalFlags.Context, globalFlags.Provider, args, toComplete, globalFlags.Owner, log.Default)
	}); err != nil {
		return err
	}

	if err := rootCmd.RegisterFlagCompletionFunc("context", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return GetContextSuggestions(rootCmd, globalFlags.Context, globalFlags.Provider, args, toComplete, globalFlags.Owner, log.Default)
	}); err != nil {
		return err
	}

	return nil
}

func GetWorkspaceSuggestions(rootCmd *cobra.Command, context, provider string, args []string, toComplete string, owner platform.OwnerFilter, logger log.Logger) ([]string, cobra.ShellCompDirective) {
	devSpaceConfig, err := config.LoadConfig(context, provider)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	workspaces, err := workspace.List(rootCmd.Context(), devSpaceConfig, false, owner, logger)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var suggestions []string
	for _, ws := range workspaces {
		if strings.HasPrefix(ws.ID, toComplete) {
			suggestions = append(suggestions, ws.ID)
		}
	}
	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

func GetProviderSuggestions(rootCmd *cobra.Command, context, provider string, args []string, toComplete string, owner platform.OwnerFilter, logger log.Logger) ([]string, cobra.ShellCompDirective) {
	devSpaceConfig, err := config.LoadConfig(context, provider)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	providers, err := workspace.LoadAllProviders(devSpaceConfig, log.Default.ErrorStreamOnly())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var suggestions []string
	for _, provider := range providers {
		if strings.HasPrefix(provider.Config.Name, toComplete) {
			suggestions = append(suggestions, provider.Config.Name)
		}
	}
	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

func GetContextSuggestions(rootCmd *cobra.Command, context, provider string, args []string, toComplete string, owner platform.OwnerFilter, logger log.Logger) ([]string, cobra.ShellCompDirective) {
	devSpaceConfig, err := config.LoadConfig(context, provider)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var suggestions []string
	for contextName, _ := range devSpaceConfig.Contexts {
		if strings.HasPrefix(contextName, toComplete) {
			suggestions = append(suggestions, contextName)
		}
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}
