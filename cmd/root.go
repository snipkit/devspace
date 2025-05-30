package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"dev.khulnasoft.com/cmd/agent"
	"dev.khulnasoft.com/cmd/completion"
	"dev.khulnasoft.com/cmd/context"
	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/cmd/helper"
	"dev.khulnasoft.com/cmd/ide"
	"dev.khulnasoft.com/cmd/machine"
	"dev.khulnasoft.com/cmd/pro"
	"dev.khulnasoft.com/cmd/provider"
	"dev.khulnasoft.com/cmd/use"
	"dev.khulnasoft.com/pkg/client/clientimplementation"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/telemetry"
	log2 "dev.khulnasoft.com/log"
	"dev.khulnasoft.com/log/terminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var globalFlags *flags.GlobalFlags

// NewRootCmd returns a new root command
func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "devspace",
		Short:         "DevSpace",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
			if globalFlags.LogOutput == "json" {
				log2.Default.SetFormat(log2.JSONFormat)
			} else if globalFlags.LogOutput == "raw" {
				log2.Default.SetFormat(log2.RawFormat)
			} else if globalFlags.LogOutput != "plain" {
				return fmt.Errorf("unrecognized log format %s, needs to be either plain or json", globalFlags.LogOutput)
			}

			if globalFlags.Silent {
				log2.Default.SetLevel(logrus.FatalLevel)
			} else if globalFlags.Debug {
				log2.Default.SetLevel(logrus.DebugLevel)
			} else if os.Getenv(clientimplementation.DevSpaceDebug) == "true" {
				log2.Default.SetLevel(logrus.DebugLevel)
			}

			if globalFlags.DevSpaceHome != "" {
				_ = os.Setenv(config.DEVSPACE_HOME, globalFlags.DevSpaceHome)
			}

			devSpaceConfig, err := config.LoadConfig(globalFlags.Context, globalFlags.Provider)
			if err == nil {
				telemetry.StartCLI(devSpaceConfig, cobraCmd)
			}

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if globalFlags.DevSpaceHome != "" {
				_ = os.Unsetenv(config.DEVSPACE_HOME)
			}

			return nil
		},
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// build the root command
	rootCmd := BuildRoot()

	// execute command
	err := rootCmd.Execute()
	telemetry.CollectorCLI.RecordCLI(err)
	telemetry.CollectorCLI.Flush()
	if err != nil {
		//nolint:all
		if sshExitErr, ok := err.(*ssh.ExitError); ok {
			os.Exit(sshExitErr.ExitStatus())
		}

		//nolint:all
		if execExitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(execExitErr.ExitCode())
		}

		if globalFlags.Debug {
			log2.Default.Fatalf("%+v", err)
		} else {
			if rootCmd.Annotations == nil || rootCmd.Annotations[agent.AgentExecutedAnnotation] != "true" {
				if terminal.IsTerminalIn {
					log2.Default.Error("Try using the --debug flag to see a more verbose output")
				} else if os.Getenv(telemetry.UIEnvVar) == "true" {
					log2.Default.Error("Try enabling Debug mode under Settings to see a more verbose output")
				}
			}
			log2.Default.Fatal(err)
		}
	}
}

// BuildRoot creates a new root command from the
func BuildRoot() *cobra.Command {
	rootCmd := NewRootCmd()
	persistentFlags := rootCmd.PersistentFlags()
	globalFlags = flags.SetGlobalFlags(persistentFlags)
	_ = completion.RegisterFlagCompletionFuns(rootCmd, globalFlags)

	rootCmd.AddCommand(agent.NewAgentCmd(globalFlags))
	rootCmd.AddCommand(provider.NewProviderCmd(globalFlags))
	rootCmd.AddCommand(use.NewUseCmd(globalFlags))
	rootCmd.AddCommand(helper.NewHelperCmd(globalFlags))
	rootCmd.AddCommand(ide.NewIDECmd(globalFlags))
	rootCmd.AddCommand(machine.NewMachineCmd(globalFlags))
	rootCmd.AddCommand(context.NewContextCmd(globalFlags))
	rootCmd.AddCommand(pro.NewProCmd(globalFlags, log2.Default))
	rootCmd.AddCommand(NewUpCmd(globalFlags))
	rootCmd.AddCommand(NewDeleteCmd(globalFlags))
	rootCmd.AddCommand(NewSSHCmd(globalFlags))
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewStopCmd(globalFlags))
	rootCmd.AddCommand(NewListCmd(globalFlags))
	rootCmd.AddCommand(NewStatusCmd(globalFlags))
	rootCmd.AddCommand(NewBuildCmd(globalFlags))
	rootCmd.AddCommand(NewLogsDaemonCmd(globalFlags))
	rootCmd.AddCommand(NewExportCmd(globalFlags))
	rootCmd.AddCommand(NewImportCmd(globalFlags))
	rootCmd.AddCommand(NewLogsCmd(globalFlags))
	rootCmd.AddCommand(NewUpgradeCmd())
	rootCmd.AddCommand(NewTroubleshootCmd(globalFlags))
	rootCmd.AddCommand(NewPingCmd(globalFlags))
	return rootCmd
}
