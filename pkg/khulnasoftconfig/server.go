package khulnasoftconfig

import (
	"encoding/json"
	"os"
	"path/filepath"

	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/pkg/provider"
)

const (
	KhulnasoftPlatformConfigFileName = "khulnasoft-config.json" // TODO: move somewhere else, replace hardoced strings with usage of this const
)

type KhulnasoftConfigRequest struct {
	// Deprecated. Do not use anymore
	Context string
	// Deprecated. Do not use anymore
	Provider string
}

type KhulnasoftConfigResponse struct {
	KhulnasoftConfig *client.Config
}

func Read(request *KhulnasoftConfigRequest) (*KhulnasoftConfigResponse, error) {
	khulnasoftConfig, err := readConfig(request.Context, request.Provider)
	if err != nil {
		return nil, err
	}

	return &KhulnasoftConfigResponse{KhulnasoftConfig: khulnasoftConfig}, nil
}

func ReadFromWorkspace(workspace *provider.Workspace) (*KhulnasoftConfigResponse, error) {
	khulnasoftConfig, err := readConfig(workspace.Context, workspace.Provider.Name)
	if err != nil {
		return nil, err
	}

	return &KhulnasoftConfigResponse{KhulnasoftConfig: khulnasoftConfig}, nil
}

func readConfig(contextName string, providerName string) (*client.Config, error) {
	providerDir, err := provider.GetProviderDir(contextName, providerName)
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(providerDir, KhulnasoftPlatformConfigFileName)

	// Check if given context and provider have Khulnasoft Platform configuration
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If not just return empty response
		return &client.Config{}, nil
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	khulnasoftConfig := &client.Config{}
	err = json.Unmarshal(content, khulnasoftConfig)
	if err != nil {
		return nil, err
	}

	return khulnasoftConfig, nil
}
