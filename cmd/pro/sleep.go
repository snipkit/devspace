package pro

import (
	"context"
	"fmt"
	"strconv"
	"time"

	clusterv1 "dev.khulnasoft.com/agentapi/pkg/apis/khulnasoft/cluster/v1"
	storagev1 "dev.khulnasoft.com/api/pkg/apis/storage/v1"
	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/config"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/platform/project"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// SleepCmd holds the cmd flags
type SleepCmd struct {
	*flags.GlobalFlags
	Log log.Logger

	Project       string
	Host          string
	ForceDuration int64
}

// NewSleepCmd creates a new command
func NewSleepCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &SleepCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Use:   "sleep",
		Short: "Put a workspace to sleep",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			log.Default.SetFormat(log.TextFormat)

			return cmd.Run(cobraCmd.Context(), args)
		},
	}

	c.Flags().StringVar(&cmd.Project, "project", "", "The project to use")
	c.Flags().Int64Var(&cmd.ForceDuration, "prevent-wakeup", -1, "The amount of seconds this workspace should sleep until it can be woken up again (use 0 for infinite sleeping). During this time the space can only be woken up by `devspace pro wakeup`, manually deleting the annotation on the namespace or through the UI")
	_ = c.MarkFlagRequired("project")
	c.Flags().StringVar(&cmd.Host, "host", "", "The pro instance to use")
	_ = c.MarkFlagRequired("host")

	return c
}

func (cmd *SleepCmd) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a workspace name")
	}
	targetWorkspace := args[0]

	devSpaceConfig, err := config.LoadConfig(cmd.Context, "")
	if err != nil {
		return err
	}

	baseClient, err := platform.InitClientFromHost(ctx, devSpaceConfig, cmd.Host, cmd.Log)
	if err != nil {
		return err
	}

	workspaceInstance, err := platform.FindInstanceByName(ctx, baseClient, targetWorkspace, cmd.Project)
	if err != nil {
		return err
	}

	managementClient, err := baseClient.Management()
	if err != nil {
		return err
	}

	// create a deep copy of the workspace instance
	oldWorkspaceInstance := workspaceInstance.DeepCopy()
	oldWorkspaceInstance.Status = workspaceInstance.Status
	oldWorkspaceInstance.ObjectMeta = workspaceInstance.ObjectMeta
	oldWorkspaceInstance.TypeMeta = workspaceInstance.TypeMeta

	// create a patch from the old workspace instance
	patch := ctrlclient.MergeFrom(oldWorkspaceInstance)
	if workspaceInstance.Annotations == nil {
		workspaceInstance.Annotations = map[string]string{}
	}
	workspaceInstance.Annotations[clusterv1.SleepModeForceAnnotation] = "true"
	if cmd.ForceDuration >= 0 {
		workspaceInstance.Annotations[clusterv1.SleepModeForceDurationAnnotation] = strconv.FormatInt(cmd.ForceDuration, 10)
	}
	patchData, err := patch.Data(workspaceInstance)
	if err != nil {
		return fmt.Errorf("create patch: %w", err)
	}

	_, err = managementClient.Khulnasoft().ManagementV1().DevSpaceWorkspaceInstances(project.ProjectNamespace(cmd.Project)).Patch(ctx, workspaceInstance.Name, patch.Type(), patchData, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	// wait for sleeping
	cmd.Log.Info("Wait until workspace is sleeping...")
	err = wait.PollUntilContextTimeout(ctx, time.Second, platform.Timeout(), false, func(ctx context.Context) (done bool, err error) {
		workspaceInstance, err := managementClient.Khulnasoft().ManagementV1().DevSpaceWorkspaceInstances(project.ProjectNamespace(cmd.Project)).Get(ctx, workspaceInstance.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		return workspaceInstance.Status.Phase == storagev1.InstanceSleeping, nil
	})
	if err != nil {
		return fmt.Errorf("error waiting for workspace to start sleeping: %w", err)
	}

	cmd.Log.Donef("Successfully put workspace %s to sleep", workspaceInstance.Name)
	return nil
}
