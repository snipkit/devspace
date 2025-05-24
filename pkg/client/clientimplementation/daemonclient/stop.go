package daemonclient

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	managementv1 "dev.khulnasoft.com/api/pkg/apis/management/v1"
	clientpkg "dev.khulnasoft.com/pkg/client"
	"dev.khulnasoft.com/pkg/platform"
)

func (c *client) Stop(ctx context.Context, opt clientpkg.StopOptions) error {
	c.m.Lock()
	defer c.m.Unlock()

	baseClient, err := c.initPlatformClient(ctx)
	if err != nil {
		return err
	}
	workspace, err := platform.FindInstance(ctx, baseClient, c.workspace.UID)
	if err != nil {
		return err
	} else if workspace == nil {
		return fmt.Errorf("couldn't find workspace")
	}

	managementClient, err := baseClient.Management()
	if err != nil {
		return err
	}

	rawOptions, _ := json.Marshal(opt)
	retStop, err := managementClient.Khulnasoft().ManagementV1().DevSpaceWorkspaceInstances(workspace.Namespace).Stop(ctx, workspace.Name, &managementv1.DevSpaceWorkspaceInstanceStop{
		Spec: managementv1.DevSpaceWorkspaceInstanceStopSpec{
			Options: string(rawOptions),
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error stopping workspace: %w", err)
	} else if retStop.Status.TaskID == "" {
		return fmt.Errorf("no stop task id returned from server")
	}

	_, err = observeTask(ctx, managementClient, workspace, retStop.Status.TaskID, c.log)
	if err != nil {
		return fmt.Errorf("stop: %w", err)
	}

	return nil
}
