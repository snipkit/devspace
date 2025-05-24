package client

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config defines the client config structure
type Config struct {
	metav1.TypeMeta `json:",inline"`

	// host is the http endpoint of how to access khulnasoft
	// +optional
	Host string `json:"host,omitempty"`

	// LastInstallContext is the last install context
	// +optional
	LastInstallContext string `json:"lastInstallContext,omitempty"`

	// insecure specifies if the khulnasoft instance is insecure
	// +optional
	Insecure bool `json:"insecure,omitempty"`

	// access key is the access key for the given khulnasoft host
	// +optional
	AccessKey string `json:"accesskey,omitempty"`

	// virtual cluster access key is the access key for the given khulnasoft host to create virtual clusters
	// +optional
	VirtualClusterAccessKey string `json:"virtualClusterAccessKey,omitempty"`

	// map of cached certificates for "access point" mode virtual clusters
	// +optional
	VirtualClusterAccessPointCertificates map[string]VirtualClusterCertificatesEntry
}

type VirtualClusterCertificatesEntry struct {
	CertificateData string
	KeyData         string
	LastRequested   metav1.Time
	ExpirationTime  time.Time
}

// NewConfig creates a new config
func NewConfig() *Config {
	return &Config{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Config",
			APIVersion: "storage.khulnasoft.com/v1",
		},
	}
}
