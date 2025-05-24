package container

import (
	"fmt"

	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/pkg/credentials"
	"dev.khulnasoft.com/pkg/khulnasoftconfig"
	"dev.khulnasoft.com/log"

	"github.com/spf13/cobra"
)

type SetupKhulnasoftPlatformAccessCmd struct {
	*flags.GlobalFlags

	Context  string
	Provider string
	Port     int
}

// NewSetupKhulnasoftPlatformAccessCmd creates a new setup-khulnasoft-platform-access command
// This agent command can be used to inject khulnasoft platform configuration from local machine to workspace.
func NewSetupKhulnasoftPlatformAccessCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &SetupKhulnasoftPlatformAccessCmd{
		GlobalFlags: flags,
	}

	setupKhulnasoftPlatformAccessCmd := &cobra.Command{
		Use:   "setup-khulnasoft-platform-access",
		Short: "used to setup Khulnasoft Platform access",
		RunE:  cmd.Run,
	}

	setupKhulnasoftPlatformAccessCmd.Flags().StringVar(&cmd.Context, "context", "", "context to use")
	_ = setupKhulnasoftPlatformAccessCmd.Flags().MarkDeprecated("context", "Information should be provided by services server, don't use this flag anymore")

	setupKhulnasoftPlatformAccessCmd.Flags().StringVar(&cmd.Provider, "provider", "", "provider to use")
	_ = setupKhulnasoftPlatformAccessCmd.Flags().MarkDeprecated("provider", "Information should be provided by services server, don't use this flag anymore")

	setupKhulnasoftPlatformAccessCmd.Flags().IntVar(&cmd.Port, "port", 0, "If specified, will use the given port")
	_ = setupKhulnasoftPlatformAccessCmd.Flags().MarkDeprecated("port", "")

	return setupKhulnasoftPlatformAccessCmd
}

// Run executes main command logic.
// It fetches Khulnasoft Platform credentials from credentials server and sets it up inside the workspace.
func (c *SetupKhulnasoftPlatformAccessCmd) Run(_ *cobra.Command, args []string) error {
	logger := log.Default.ErrorStreamOnly()

	port, err := credentials.GetPort()
	if err != nil {
		return fmt.Errorf("get port: %w", err)
	}
	// backwards compatibility, remove in future release
	if c.Port > 0 {
		port = c.Port
	}

	khulnasoftConfig, err := khulnasoftconfig.GetKhulnasoftConfig(c.Context, c.Provider, port, logger)
	if err != nil {
		return err
	}

	if khulnasoftConfig == nil {
		logger.Debug("Got empty khulnasoft config response, Khulnasoft Platform access won't be set up.")
		return nil
	}

	err = khulnasoftconfig.AuthDevspaceCliToPlatform(khulnasoftConfig, logger)
	if err != nil {
		// log error but don't return to allow other CLIs to install as well
		logger.Warnf("unable to authenticate devspace cli: %w", err)
	}

	err = khulnasoftconfig.AuthVClusterCliToPlatform(khulnasoftConfig, logger)
	if err != nil {
		// log error but don't return to allow other CLIs to install as well
		logger.Warnf("unable to authenticate vcluster cli: %w", err)
	}

	return nil
}
