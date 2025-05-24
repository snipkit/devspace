package khulnasoftconfig

import (
	"bytes"
	"encoding/json"
	"fmt"

	"dev.khulnasoft.com/pkg/credentials"
	"dev.khulnasoft.com/pkg/platform/client"
	"dev.khulnasoft.com/log"
)

func GetKhulnasoftConfig(context, provider string, port int, logger log.Logger) (*client.Config, error) {
	request := &KhulnasoftConfigRequest{
		Context:  context,
		Provider: provider,
	}

	rawJson, err := json.Marshal(request)
	if err != nil {
		logger.Errorf("Error parsing request: %w", err)
		return nil, err
	}

	configResponse := &KhulnasoftConfigResponse{}
	out, err := credentials.PostWithRetry(port, "khulnasoft-platform-credentials", bytes.NewReader(rawJson), logger)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(out, configResponse)
	if err != nil {
		return nil, fmt.Errorf("decode khulnasoft config %s: %w", string(out), err)
	}

	return configResponse.KhulnasoftConfig, nil
}
