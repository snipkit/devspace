package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"dev.khulnasoft.com/cmd/pro/flags"
	devspacehttp "dev.khulnasoft.com/pkg/http"
	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/log"
	"github.com/spf13/cobra"
)

// HealthCmd holds the cmd flags
type HealthCmd struct {
	*flags.GlobalFlags

	Log log.Logger
}

// NewHealthCmd creates a new command
func NewHealthCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &HealthCmd{
		GlobalFlags: globalFlags,
		Log:         log.GetInstance(),
	}
	c := &cobra.Command{
		Use:    "health",
		Short:  "Check platform health",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run(cobraCmd.Context(), os.Stdin, os.Stdout, os.Stderr)
		},
	}

	return c
}

func (cmd *HealthCmd) Run(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	baseClient, err := client.InitClientFromPath(ctx, cmd.Config)
	if err != nil {
		return err
	}

	config := baseClient.Config()
	u, err := url.Parse(fmt.Sprintf("%s/healthz", config.Host))
	if err != nil {
		return err
	}
	res, err := devspacehttp.GetHTTPClient().Get(u.String())
	if err != nil {
		return err
	}
	healthCheck := HealthCheck{
		Healthy: res.StatusCode == http.StatusOK,
	}

	out, err := json.Marshal(healthCheck)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}

type HealthCheck struct {
	Healthy bool `json:"healthy,omitempty"`
}
