package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	managementv1 "dev.khulnasoft.com/api/pkg/apis/management/v1"
	"dev.khulnasoft.com/cmd/completion"
	"dev.khulnasoft.com/cmd/flags"
	"dev.khulnasoft.com/cmd/provider"
	"dev.khulnasoft.com/pkg/client"
	"dev.khulnasoft.com/pkg/config"
	daemon "dev.khulnasoft.com/pkg/daemon/platform"
	"dev.khulnasoft.com/pkg/platform"
	pkgprovider "dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/pkg/version"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TroubleshootCmd struct {
	*flags.GlobalFlags
}

func NewTroubleshootCmd(flags *flags.GlobalFlags) *cobra.Command {
	cmd := &TroubleshootCmd{
		GlobalFlags: flags,
	}
	troubleshootCmd := &cobra.Command{
		Use:   "troubleshoot [workspace-path|workspace-name]",
		Short: "Prints the workspaces troubleshooting information",
		Run: func(cobraCmd *cobra.Command, args []string) {
			cmd.Run(cobraCmd.Context(), args)
		},
		ValidArgsFunction: func(rootCmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completion.GetWorkspaceSuggestions(rootCmd, cmd.Context, cmd.Provider, args, toComplete, cmd.Owner, log.Default)
		},
		Hidden: true,
	}

	return troubleshootCmd
}

func (cmd *TroubleshootCmd) Run(ctx context.Context, args []string) {
	// (ThomasK33): We're creating an anonymous struct here, so that we group
	// everything and then we can serialize it in one call.
	var info struct {
		CLIVersion            string
		Config                *config.Config
		Providers             map[string]provider.ProviderWithDefault
		DevSpaceProInstances    []DevSpaceProInstance
		Workspace             *pkgprovider.Workspace
		WorkspaceStatus       client.Status
		WorkspaceTroubleshoot *managementv1.DevSpaceWorkspaceInstanceTroubleshoot
		DaemonStatus          *daemon.Status

		Errors []PrintableError `json:",omitempty"`
	}
	info.CLIVersion = version.GetVersion()

	// (ThomasK33): We are defering the printing here, as we want to make sure
	// that we will always print, even in the case of a panic.
	defer func() {
		out, err := json.MarshalIndent(info, "", "  ")
		if err == nil {
			fmt.Print(string(out))
		} else {
			fmt.Print(err)
			fmt.Print(info)
		}
	}()

	// NOTE(ThomasK33): Since this is a troubleshooting command, we want to
	// collect as many relevant information as possible.
	// For this reason we may not return with an error early.
	// We are fine with a partially filled TrbouelshootInfo struct, as this
	// already provides us with more information then before.
	var err error
	info.Config, err = config.LoadConfig(cmd.Context, cmd.Provider)
	if err != nil {
		info.Errors = append(info.Errors, PrintableError{fmt.Errorf("load config: %w", err)})
		// (ThomasK33): It's fine to return early here, as without the devspace config
		// we cannot do any further troubleshooting.
		return
	}

	logger := log.Default.ErrorStreamOnly()
	info.Providers, err = collectProviders(info.Config, logger)
	if err != nil {
		info.Errors = append(info.Errors, PrintableError{fmt.Errorf("collect providers: %w", err)})
	}

	info.DevSpaceProInstances, err = collectPlatformInfo(info.Config, logger)
	if err != nil {
		info.Errors = append(info.Errors, PrintableError{fmt.Errorf("collect platform info: %w", err)})
	}

	workspaceClient, err := workspace.Get(ctx, info.Config, args, false, cmd.Owner, false, logger)
	if err == nil {
		info.Workspace = workspaceClient.WorkspaceConfig()
		info.WorkspaceStatus, err = workspaceClient.Status(ctx, client.StatusOptions{})
		if err != nil {
			info.Errors = append(info.Errors, PrintableError{fmt.Errorf("workspace status: %w", err)})
		}

		if info.Workspace.Pro != nil {
			// (ThomasK33): As there can be multiple pro instances configured
			// we want to iterate over all and find the host that this workspace belongs to.
			var proInstance DevSpaceProInstance

			for _, instance := range info.DevSpaceProInstances {
				if instance.ProviderName == info.Workspace.Provider.Name {
					proInstance = instance
					break
				}
			}

			if proInstance.ProviderName != "" {
				info.WorkspaceTroubleshoot, err = collectProWorkspaceInfo(
					ctx,
					info.Config,
					proInstance.Host,
					logger,
					info.Workspace.UID,
					info.Workspace.Pro.Project,
				)
				if err != nil {
					info.Errors = append(info.Errors, PrintableError{fmt.Errorf("collect pro workspace info: %w", err)})
				}
			}
		}
	} else {
		info.Errors = append(info.Errors, PrintableError{fmt.Errorf("get workspace: %w", err)})
	}

	daemonClient, ok := workspaceClient.(client.DaemonClient)
	if ok {
		status, err := daemon.NewLocalClient(daemonClient.Provider()).Status(ctx, true)
		if err != nil {
			info.Errors = append(info.Errors, PrintableError{fmt.Errorf("get daemon status: %w", err)})
		} else {
			info.DaemonStatus = &status
		}
	}
}

// collectProWorkspaceInfo collects troubleshooting information for a DevSpace Pro instance.
// It initializes a client from the host, finds the workspace instance in the project, and retrieves
// troubleshooting information using the management client.
func collectProWorkspaceInfo(
	ctx context.Context,
	devSpaceConfig *config.Config,
	host string,
	logger log.Logger,
	workspaceUID string,
	project string,
) (*managementv1.DevSpaceWorkspaceInstanceTroubleshoot, error) {
	baseClient, err := platform.InitClientFromHost(ctx, devSpaceConfig, host, logger)
	if err != nil {
		return nil, fmt.Errorf("init client from host: %w", err)
	}

	workspace, err := platform.FindInstanceInProject(ctx, baseClient, workspaceUID, project)
	if err != nil {
		return nil, err
	} else if workspace == nil {
		return nil, fmt.Errorf("couldn't find workspace")
	}

	managementClient, err := baseClient.Management()
	if err != nil {
		return nil, fmt.Errorf("management: %w", err)
	}

	troubleshoot, err := managementClient.
		Loft().
		ManagementV1().
		DevSpaceWorkspaceInstances(workspace.Namespace).
		Troubleshoot(ctx, workspace.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("troubleshoot: %w", err)
	}

	return troubleshoot, nil
}

// collectProviders collects and configures providers based on the given devSpaceConfig.
// It returns a map of providers with their default settings and an error if any occurs.
func collectProviders(devSpaceConfig *config.Config, logger log.Logger) (map[string]provider.ProviderWithDefault, error) {
	providers, err := workspace.LoadAllProviders(devSpaceConfig, logger)
	if err != nil {
		return nil, err
	}

	configuredProviders := devSpaceConfig.Current().Providers
	if configuredProviders == nil {
		configuredProviders = map[string]*config.ProviderConfig{}
	}

	retMap := map[string]provider.ProviderWithDefault{}
	for k, entry := range providers {
		if configuredProviders[entry.Config.Name] == nil {
			continue
		}

		srcOptions := provider.MergeDynamicOptions(entry.Config.Options, configuredProviders[entry.Config.Name].DynamicOptions)
		entry.Config.Options = srcOptions
		retMap[k] = provider.ProviderWithDefault{
			ProviderWithOptions: *entry,
			Default:             devSpaceConfig.Current().DefaultProvider == entry.Config.Name,
		}
	}

	return retMap, nil
}

type DevSpaceProInstance struct {
	Host         string
	ProviderName string
	Version      string
}

// collectPlatformInfo collects information about all platform instances in a given devSpaceConfig.
// It iterates over the pro instances, retrieves their versions, and appends them to the ProInstance slice.
// Any errors encountered during this process are combined and returned along with the ProInstance slice.
// This means that even when an error value is returned, the pro instance slice will contain valid values.
func collectPlatformInfo(devSpaceConfig *config.Config, logger log.Logger) ([]DevSpaceProInstance, error) {
	proInstanceList, err := workspace.ListProInstances(devSpaceConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("list pro instances: %w", err)
	}

	var proInstances []DevSpaceProInstance
	var combinedErrs error

	for _, proInstance := range proInstanceList {
		version, err := platform.GetProInstanceDevSpaceVersion(&pkgprovider.ProInstance{Host: proInstance.Host})
		combinedErrs = errors.Join(combinedErrs, err)
		proInstances = append(proInstances, DevSpaceProInstance{
			Host:         proInstance.Host,
			ProviderName: proInstance.Provider,
			Version:      version,
		})
	}

	return proInstances, combinedErrs
}

// (ThomasK33): Little type embedding here, so that we can
// serialize the error strings when invoking json.Marshal.
type PrintableError struct{ error }

func (p PrintableError) MarshalJSON() ([]byte, error) { return json.Marshal(p.Error()) }
