package pro

import (
	"context"
	"fmt"
	"strconv"

	proflags "dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/cmd/pro/provider/list"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/options"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/pkg/platform/parameters"
	"dev.khulnasoft.com/pkg/platform/project"
	provider2 "dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/pkg/random"
	"dev.khulnasoft.com/pkg/workspace"
	"dev.khulnasoft.com/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	managementv1 "dev.khulnasoft.com/api/pkg/apis/management/v1"
	storagev1 "dev.khulnasoft.com/api/pkg/apis/storage/v1"
)

type ImportCmd struct {
	*proflags.GlobalFlags

	WorkspaceId      string
	WorkspaceUid     string
	WorkspaceProject string

	Own bool
	log log.Logger
}

// NewImportCmd creates a new command
func NewImportCmd(globalFlags *proflags.GlobalFlags) *cobra.Command {
	logger := log.GetInstance()
	cmd := &ImportCmd{
		GlobalFlags: globalFlags,
		log:         logger,
	}

	importCmd := &cobra.Command{
		Use:   "import-workspace",
		Short: "Imports a workspace",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run(cobraCmd.Context(), args)
		},
	}

	importCmd.Flags().StringVar(&cmd.WorkspaceId, "workspace-id", "", "ID of a workspace to import")
	importCmd.Flags().StringVar(&cmd.WorkspaceUid, "workspace-uid", "", "UID of a workspace to import")
	importCmd.Flags().StringVar(&cmd.WorkspaceProject, "workspace-project", "", "Project of the workspace to import")
	importCmd.Flags().BoolVar(&cmd.Own, "own", false, "If true, will behave as if workspace was not imported")
	_ = importCmd.MarkFlagRequired("workspace-uid")
	return importCmd
}

func (cmd *ImportCmd) Run(ctx context.Context, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: devspace pro import-workspace <devspace-pro-host>")
	}

	devSpaceProHost := args[0]
	devSpaceConfig, err := config.LoadConfig(cmd.Context, "")
	if err != nil {
		return err
	}

	// set uid as id
	if cmd.WorkspaceId == "" {
		cmd.WorkspaceId = cmd.WorkspaceUid
	}

	// check if workspace already exists
	if provider2.WorkspaceExists(devSpaceConfig.DefaultContext, cmd.WorkspaceId) {
		workspaceConfig, err := provider2.LoadWorkspaceConfig(devSpaceConfig.DefaultContext, cmd.WorkspaceId)
		if err != nil {
			return fmt.Errorf("load workspace: %w", err)
		} else if workspaceConfig.UID == cmd.WorkspaceUid {
			cmd.log.Infof("Workspace %s already imported", cmd.WorkspaceId)
			return nil
		}

		newWorkspaceId := cmd.WorkspaceId + "-" + random.String(5)
		if provider2.WorkspaceExists(devSpaceConfig.DefaultContext, newWorkspaceId) {
			return fmt.Errorf("workspace %s already exists", cmd.WorkspaceId)
		}

		cmd.log.Infof("Workspace %s already exists, will use name %s instead", cmd.WorkspaceId, newWorkspaceId)
		cmd.WorkspaceId = newWorkspaceId
	}

	provider, err := workspace.ProviderFromHost(ctx, devSpaceConfig, devSpaceProHost, cmd.log)
	if err != nil {
		return fmt.Errorf("resolve provider: %w", err)
	}

	baseClient, err := platform.InitClientFromProvider(ctx, devSpaceConfig, provider.Name, cmd.log)
	if err != nil {
		return fmt.Errorf("base client: %w", err)
	}
	instance, err := platform.FindInstanceInProject(ctx, baseClient, cmd.WorkspaceUid, cmd.WorkspaceProject)
	if err != nil {
		return fmt.Errorf("find workspace instance: %w", err)
	}
	if instance == nil {
		return fmt.Errorf("workspace instance with UID %s not found", cmd.WorkspaceUid)
	}

	// old pro provider
	if !provider.HasHealthCheck() {
		instanceOpts, err := resolveInstanceOptions(ctx, instance, baseClient)
		if err != nil {
			return fmt.Errorf("resolve instance options: %w", err)
		}

		err = cmd.writeWorkspaceDefinition(devSpaceConfig, provider, instanceOpts, instance)
		if err != nil {
			return errors.Wrap(err, "prepare workspace to import definition")
		}
		cmd.log.Infof("Successfully imported workspace %s", cmd.WorkspaceId)
		return nil
	}

	// new pro provider
	err = cmd.writeNewWorkspaceDefinition(devSpaceConfig, instance, provider.Name)
	if err != nil {
		return errors.Wrap(err, "prepare workspace to import definition")
	}
	cmd.log.Infof("Successfully imported workspace %s", cmd.WorkspaceId)

	return nil
}

func (cmd *ImportCmd) writeNewWorkspaceDefinition(devSpaceConfig *config.Config, instance *managementv1.DevSpaceWorkspaceInstance, providerName string) error {
	workspaceObj := &provider2.Workspace{
		ID:       cmd.WorkspaceId,
		UID:      cmd.WorkspaceUid,
		Provider: provider2.WorkspaceProviderConfig{Name: providerName},
		Context:  devSpaceConfig.DefaultContext,
		Imported: !cmd.Own,
		Pro: &provider2.ProMetadata{
			InstanceName: instance.GetName(),
			Project:      project.ProjectFromNamespace(instance.Namespace),
			DisplayName:  instance.Spec.DisplayName,
		},
	}

	return provider2.SaveWorkspaceConfig(workspaceObj)
}

func (cmd *ImportCmd) writeWorkspaceDefinition(devSpaceConfig *config.Config, provider *provider2.ProviderConfig, instanceOpts map[string]string, instance *managementv1.DevSpaceWorkspaceInstance) error {
	workspaceObj := &provider2.Workspace{
		ID:  cmd.WorkspaceId,
		UID: cmd.WorkspaceUid,
		Provider: provider2.WorkspaceProviderConfig{
			Name:    provider.Name,
			Options: map[string]config.OptionValue{},
		},
		Context:  devSpaceConfig.DefaultContext,
		Imported: !cmd.Own,
		Pro: &provider2.ProMetadata{
			InstanceName: instance.GetName(),
			Project:      instanceOpts[platform.ProjectEnv],
			DisplayName:  instance.Spec.DisplayName,
		},
	}

	devSpaceConfig, err := options.ResolveOptions(context.Background(), devSpaceConfig, provider, instanceOpts, false, false, nil, cmd.log)
	if err != nil {
		return fmt.Errorf("resolve options: %w", err)
	}
	if devSpaceConfig.Current() == nil || devSpaceConfig.Current().Providers[provider.Name] == nil {
		return fmt.Errorf("unable to resolve provider config for provider %s", provider.Name)
	}
	workspaceObj.Provider.Options = devSpaceConfig.Current().Providers[provider.Name].Options

	err = provider2.SaveWorkspaceConfig(workspaceObj)
	if err != nil {
		return err
	}

	return nil
}

func resolveInstanceOptions(ctx context.Context, instance *managementv1.DevSpaceWorkspaceInstance, baseClient client.Client) (map[string]string, error) {
	opts := map[string]string{}
	projectName := project.ProjectFromNamespace(instance.Namespace)

	opts[platform.ProjectEnv] = projectName
	if instance.Spec.TemplateRef == nil {
		return opts, nil
	}
	//nolint:all
	if instance.Spec.RunnerRef.Runner != "" {
		opts[platform.RunnerEnv] = instance.Spec.RunnerRef.Runner //nolint:all
	}
	opts[platform.TemplateOptionEnv] = instance.Spec.TemplateRef.Name

	if instance.Spec.TemplateRef.Version != "" {
		opts[platform.TemplateVersionOptionEnv] = instance.Spec.TemplateRef.Version
	}

	if instance.Spec.Parameters == "" {
		return opts, nil
	}
	managementClient, err := baseClient.Management()
	if err != nil {
		return nil, fmt.Errorf("get management client: %w", err)
	}
	template, err := list.FindTemplate(ctx, managementClient, projectName, instance.Spec.TemplateRef.Name)
	if err != nil {
		return nil, fmt.Errorf("find template: %w", err)
	}
	templateParameters := template.Spec.Parameters
	if len(template.Spec.Versions) > 0 {
		templateParameters, err = list.GetTemplateParameters(template, instance.Spec.TemplateRef.Version)
		if err != nil {
			return nil, fmt.Errorf("get template parameters: %w", err)
		}
	}
	err = fillParameterOptions(opts, templateParameters, instance.Spec.Parameters)
	if err != nil {
		return nil, fmt.Errorf("fill parameter options: %w", err)
	}

	return opts, nil
}

func fillParameterOptions(opts map[string]string, parameterDefinitions []storagev1.AppParameter, instanceParameters string) error {
	parametersMap := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(instanceParameters), &parametersMap)
	if err != nil {
		return fmt.Errorf("unmarshal parameters: %w", err)
	}

	for _, parameter := range parameterDefinitions {
		val := parameters.GetDeepValue(parametersMap, parameter.Variable)
		var strVal string
		if val != nil {
			switch t := val.(type) {
			case string:
				strVal = t
			case int:
				strVal = strconv.Itoa(t)
			case bool:
				strVal = strconv.FormatBool(t)
			default:
				return fmt.Errorf("unrecognized type for parameter %s (%s) in file: %v", parameter.Label, parameter.Variable, t)
			}
		}

		_, err := parameters.VerifyValue(strVal, parameter)
		if err != nil {
			return err
		}

		optionName := list.VariableToEnvironmentVariable(parameter.Variable)
		opts[optionName] = strVal
	}

	return nil
}
