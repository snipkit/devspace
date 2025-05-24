package provider

import "dev.khulnasoft.com/pkg/types"

type ProInstance struct {
	// Provider is the provider name this pro instance belongs to
	Provider string `json:"provider,omitempty"`

	// Host is the Loft DevSpace Pro host to use
	Host string `json:"host,omitempty"`

	// CreationTimestamp is the timestamp when this pro instance was created
	CreationTimestamp types.Time `json:"creationTimestamp,omitempty"`
}
