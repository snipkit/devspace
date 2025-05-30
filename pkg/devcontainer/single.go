package devcontainer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"dev.khulnasoft.com/pkg/daemon/agent"
	"dev.khulnasoft.com/pkg/devcontainer/config"
	"dev.khulnasoft.com/pkg/devcontainer/metadata"
	"dev.khulnasoft.com/pkg/driver"
	provider2 "dev.khulnasoft.com/pkg/provider"
	"github.com/pkg/errors"
)

var dockerlessImage = "ghcr.io/khulnasoft-lab/dockerless:0.2.0"

const (
	DevSpaceExtraEnvVar           = "DEVSPACE"
	RemoteContainersExtraEnvVar = "REMOTE_CONTAINERS"
	WorkspaceIDExtraEnvVar      = "DEVSPACE_WORKSPACE_ID"
	WorkspaceUIDExtraEnvVar     = "DEVSPACE_WORKSPACE_UID"

	DefaultEntrypoint = `
while ! command -v /usr/local/bin/devspace >/dev/null 2>&1; do
  echo "Waiting for devspace tool..."
  sleep 1
done
exec /usr/local/bin/devspace agent container daemon
`
)

func (r *runner) runSingleContainer(
	ctx context.Context,
	parsedConfig *config.SubstitutedConfig,
	substitutionContext *config.SubstitutionContext,
	options UpOptions,
	timeout time.Duration,
) (*config.Result, error) {
	r.Log.Debugf("Starting devcontainer in single container mode...")
	containerDetails, err := r.Driver.FindDevContainer(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("find dev container: %w", err)
	}

	// does the container already exist?
	var (
		mergedConfig *config.MergedDevContainerConfig
	)

	// if options.Recreate is true, and workspace is a running container, we should not rebuild
	if options.Recreate && parsedConfig.Config.ContainerID != "" {
		return nil, fmt.Errorf("cannot recreate container not created by DevSpace")
	} else if !options.Recreate && containerDetails != nil {
		// start container if not running
		if strings.ToLower(containerDetails.State.Status) != "running" {
			err = r.Driver.StartDevContainer(ctx, r.ID)
			if err != nil {
				return nil, err
			}
		}

		// if we are working with a non-managed container, and it has set workingDir, set it as the workspaceFolder
		if parsedConfig.Config.ContainerID != "" && containerDetails.Config.WorkingDir != "" {
			substitutionContext.ContainerWorkspaceFolder = containerDetails.Config.WorkingDir
		}

		imageMetadataConfig, err := metadata.GetImageMetadataFromContainer(containerDetails, substitutionContext, r.Log)
		if err != nil {
			return nil, err
		}

		mergedConfig, err = config.MergeConfiguration(parsedConfig.Config, imageMetadataConfig.Config)
		if err != nil {
			return nil, errors.Wrap(err, "merge config")
		}

		// If driver can reprovision, rerun the devcontainer and let the driver handle follow-up steps
		if d, ok := r.Driver.(driver.ReprovisioningDriver); ok && d.CanReprovision() {
			err = r.Driver.RunDevContainer(ctx, r.ID, nil)
			if err != nil {
				return nil, errors.Wrap(err, "start dev container")
			}

			// get from build info
			containerDetails, err = r.Driver.FindDevContainer(ctx, r.ID)
			if err != nil {
				return nil, fmt.Errorf("find dev container: %w", err)
			}
		}
	} else {
		// we need to build the container
		buildInfo, err := r.build(ctx, parsedConfig, substitutionContext, provider2.BuildOptions{
			CLIOptions: provider2.CLIOptions{
				PrebuildRepositories: options.PrebuildRepositories,
				ForceDockerless:      options.ForceDockerless,
				Platform:             options.CLIOptions.Platform,
			},
			NoBuild:       options.NoBuild,
			RegistryCache: options.RegistryCache,
			ExportCache:   false,
		})
		if err != nil {
			return nil, errors.Wrap(err, "build image")
		}

		// delete container on recreation
		if options.Recreate {
			if _, ok := r.Driver.(driver.DockerDriver); ok {
				err := r.Delete(ctx)
				if err != nil {
					return nil, errors.Wrap(err, "delete devcontainer")
				}
			} else {
				err := r.Driver.StopDevContainer(ctx, r.ID)
				if err != nil {
					return nil, errors.Wrap(err, "stop devcontainer")
				}
			}
		}

		// merge configuration
		mergedConfig, err = config.MergeConfiguration(parsedConfig.Config, buildInfo.ImageMetadata.Config)
		if err != nil {
			return nil, errors.Wrap(err, "merge config")
		}

		// Inject the daemon entrypoint if platform configuration is provided.
		if options.CLIOptions.Platform.AccessKey != "" {
			r.Log.Debugf("Platform config detected, injecting DevSpace daemon entrypoint.")

			data, err := agent.GetEncodedWorkspaceDaemonConfig(options.Platform, r.WorkspaceConfig.Workspace, substitutionContext, mergedConfig)
			if err != nil {
				r.Log.Errorf("Failed to marshal daemon config: %v", err)
			} else {
				mergedConfig.ContainerEnv[config.WorkspaceDaemonConfigExtraEnvVar] = data
			}
		}

		// run dev container
		err = r.runContainer(ctx, parsedConfig, substitutionContext, mergedConfig, buildInfo)
		if err != nil {
			return nil, errors.Wrap(err, "start dev container")
		}

		// TODO: wait here a bit for correct startup?

		// get from build info
		containerDetails, err = r.Driver.FindDevContainer(ctx, r.ID)
		if err != nil {
			return nil, fmt.Errorf("find dev container: %w", err)
		}
	}

	// setup container
	return r.setupContainer(ctx, parsedConfig.Raw, containerDetails, mergedConfig, substitutionContext, timeout)
}

func (r *runner) runContainer(
	ctx context.Context,
	parsedConfig *config.SubstitutedConfig,
	substitutionContext *config.SubstitutionContext,
	mergedConfig *config.MergedDevContainerConfig,
	buildInfo *config.BuildInfo,
) error {
	var err error

	// build run options for dockerless mode
	var runOptions *driver.RunOptions
	if buildInfo.Dockerless != nil {
		runOptions, err = r.getDockerlessRunOptions(mergedConfig, substitutionContext, buildInfo)
		if err != nil {
			return fmt.Errorf("build dockerless run options: %w", err)
		}
	} else {
		// build run options
		runOptions, err = r.getRunOptions(mergedConfig, substitutionContext, buildInfo)
		if err != nil {
			return fmt.Errorf("build run options: %w", err)
		}
	}

	runOptions.Env = r.addExtraEnvVars(runOptions.Env)

	// check if docker
	dockerDriver, ok := r.Driver.(driver.DockerDriver)
	if ok {
		return dockerDriver.RunDockerDevContainer(
			ctx,
			r.ID,
			runOptions,
			parsedConfig.Config,
			mergedConfig.Init,
			r.WorkspaceConfig.Workspace.IDE.Name,
			r.WorkspaceConfig.Workspace.IDE.Options,
		)
	}

	// build run options for regular driver
	return r.Driver.RunDevContainer(ctx, r.ID, runOptions)
}

func (r *runner) getDockerlessRunOptions(
	mergedConfig *config.MergedDevContainerConfig,
	substitutionContext *config.SubstitutionContext,
	buildInfo *config.BuildInfo,
) (*driver.RunOptions, error) {
	// parse workspace mount
	workspaceMountParsed := config.ParseMount(substitutionContext.WorkspaceMount)

	// add metadata as label here
	marshalled, err := json.Marshal(buildInfo.ImageMetadata.Raw)
	if err != nil {
		return nil, errors.Wrap(err, "marshal config")
	}
	env := map[string]string{
		"DOCKERLESS":            "true",
		"DOCKERLESS_CONTEXT":    buildInfo.Dockerless.Context,
		"DOCKERLESS_DOCKERFILE": buildInfo.Dockerless.Dockerfile,
		"GODEBUG":               "http2client=0", // https://github.com/GoogleContainerTools/kaniko/issues/875
	}
	for k, v := range mergedConfig.ContainerEnv {
		env[k] = v
	}
	if buildInfo.Dockerless.Target != "" {
		env["DOCKERLESS_TARGET"] = buildInfo.Dockerless.Target
	}
	if len(buildInfo.Dockerless.BuildArgs) > 0 {
		out, err := json.Marshal(config.ObjectToList(buildInfo.Dockerless.BuildArgs))
		if err != nil {
			return nil, fmt.Errorf("marshal build args: %w", err)
		}

		env["DOCKERLESS_BUILD_ARGS"] = string(out)
	}

	image := dockerlessImage
	if r.WorkspaceConfig != nil && r.WorkspaceConfig.Agent.Dockerless.Image != "" {
		image = r.WorkspaceConfig.Agent.Dockerless.Image
	}

	// we need to add an extra mount here, because otherwise the build config might get lost
	mounts := mergedConfig.Mounts
	mounts = append(mounts, &config.Mount{
		Type:   "volume",
		Source: "dockerless-" + r.ID,
		Target: "/workspaces/.dockerless",
	})

	uid := ""
	if r.WorkspaceConfig != nil && r.WorkspaceConfig.Workspace != nil {
		uid = r.WorkspaceConfig.Workspace.UID
	}

	// build run options
	return &driver.RunOptions{
		UID:        uid,
		Image:      image,
		User:       "root",
		Entrypoint: "/.dockerless/dockerless",
		Cmd: []string{
			"start",
			"--wait",
			"--entrypoint", "/.dockerless/bin/sh",
			"--cmd", "-c",
			"--cmd", GetStartScript(mergedConfig),
			"--user", buildInfo.Dockerless.User,
		},
		Env:    env,
		CapAdd: mergedConfig.CapAdd,
		Labels: []string{
			metadata.ImageMetadataLabel + "=" + string(marshalled),
			config.UserLabel + "=" + buildInfo.Dockerless.User,
		},
		Privileged:     mergedConfig.Privileged,
		WorkspaceMount: &workspaceMountParsed,
		Mounts:         mounts,
	}, nil
}

func (r *runner) getRunOptions(
	mergedConfig *config.MergedDevContainerConfig,
	substitutionContext *config.SubstitutionContext,
	buildInfo *config.BuildInfo,
) (*driver.RunOptions, error) {
	// parse workspace mount
	workspaceMountParsed := config.ParseMount(substitutionContext.WorkspaceMount)

	// add metadata as label here
	marshalled, err := json.Marshal(buildInfo.ImageMetadata.Raw)
	if err != nil {
		return nil, errors.Wrap(err, "marshal config")
	}

	// build labels & entrypoint
	entrypoint, cmd := GetContainerEntrypointAndArgs(mergedConfig, buildInfo.ImageDetails)
	labels := []string{
		metadata.ImageMetadataLabel + "=" + string(marshalled),
		config.UserLabel + "=" + buildInfo.ImageDetails.Config.User,
	}

	user := buildInfo.ImageDetails.Config.User
	if mergedConfig.ContainerUser != "" {
		user = mergedConfig.ContainerUser
	}

	uid := ""
	if r.WorkspaceConfig != nil && r.WorkspaceConfig.Workspace != nil {
		uid = r.WorkspaceConfig.Workspace.UID
	}

	return &driver.RunOptions{
		UID:            uid,
		Image:          buildInfo.ImageName,
		User:           user,
		Entrypoint:     entrypoint,
		Cmd:            cmd,
		Env:            mergedConfig.ContainerEnv,
		CapAdd:         mergedConfig.CapAdd,
		Labels:         labels,
		Privileged:     mergedConfig.Privileged,
		WorkspaceMount: &workspaceMountParsed,
		SecurityOpt:    mergedConfig.SecurityOpt,
		Mounts:         mergedConfig.Mounts,
	}, nil
}

// add environment variables that signals that we are in a remote container
// (vscode compatibility) and specifically that we are using devspace.
func (r *runner) addExtraEnvVars(env map[string]string) map[string]string {
	if env == nil {
		env = make(map[string]string)
	}

	env[DevSpaceExtraEnvVar] = "true"
	env[RemoteContainersExtraEnvVar] = "true"
	if r.WorkspaceConfig != nil && r.WorkspaceConfig.Workspace != nil && r.WorkspaceConfig.Workspace.ID != "" {
		env[WorkspaceIDExtraEnvVar] = r.WorkspaceConfig.Workspace.ID
	}
	if r.WorkspaceConfig != nil && r.WorkspaceConfig.Workspace != nil && r.WorkspaceConfig.Workspace.UID != "" {
		env[WorkspaceUIDExtraEnvVar] = r.WorkspaceConfig.Workspace.UID
	}

	return env
}

func GetStartScript(mergedConfig *config.MergedDevContainerConfig) string {
	customEntrypoints := mergedConfig.Entrypoints
	return `echo Container started
trap "exit 0" 15
exec "$@"
` + strings.Join(customEntrypoints, "\n") + DefaultEntrypoint
}

func GetContainerEntrypointAndArgs(mergedConfig *config.MergedDevContainerConfig, imageDetails *config.ImageDetails) (string, []string) {
	cmd := []string{"-c", GetStartScript(mergedConfig), "-"} // `wait $!` allows for the `trap` to run (synchronous `sleep` would not).
	if imageDetails != nil && mergedConfig.OverrideCommand != nil && !*mergedConfig.OverrideCommand {
		cmd = append(cmd, imageDetails.Config.Entrypoint...)
		cmd = append(cmd, imageDetails.Config.Cmd...)
	}
	return "/bin/sh", cmd
}
