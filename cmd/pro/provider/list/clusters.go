package list

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	managementv1 "dev.khulnasoft.com/api/pkg/apis/management/v1"
	"dev.khulnasoft.com/cmd/pro/flags"
	"dev.khulnasoft.com/pkg/platform"
	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClustersCmd holds the cmd flags
type ClustersCmd struct {
	*flags.GlobalFlags

	log log.Logger
}

// NewClustersCmd creates a new command
func NewClustersCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &ClustersCmd{
		GlobalFlags: globalFlags,
		log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Use:   "clusters",
		Short: "Lists clusters for the provider",
		Args:  cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run(cobraCmd.Context())
		},
	}

	return c
}

func (cmd *ClustersCmd) Run(ctx context.Context) error {
	projectName := os.Getenv(platform.ProjectEnv)
	if projectName == "" {
		return fmt.Errorf("%s environment variable is empty", platform.ProjectEnv)
	}

	baseClient, err := client.InitClientFromPath(ctx, cmd.Config)
	if err != nil {
		return err
	}

	clustersList, err := Clusters(ctx, baseClient, projectName)
	if err != nil {
		return err
	}

	out, err := json.Marshal(clustersList)
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	return nil
}

func Clusters(ctx context.Context, client client.Client, projectName string) (*managementv1.ProjectClusters, error) {
	managementClient, err := client.Management()
	if err != nil {
		return nil, err
	}

	clustersList, err := managementClient.Loft().ManagementV1().Projects().ListClusters(ctx, projectName, metav1.GetOptions{})
	if err != nil {
		return clustersList, fmt.Errorf("list clusters: %w", err)
	}

	return clustersList, nil
}
