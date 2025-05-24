package kube

import (
	"fmt"

	khulnasoftclient "dev.khulnasoft.com/api/pkg/clientset/versioned"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Interface interface {
	kubernetes.Interface
	Khulnasoft() khulnasoftclient.Interface
}

func NewForConfig(c *rest.Config) (Interface, error) {
	kubeClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, fmt.Errorf("create kube client: %w", err)
	}

	khulnasoftClient, err := khulnasoftclient.NewForConfig(c)
	if err != nil {
		return nil, fmt.Errorf("create khulnasoft client: %w", err)
	}

	return &client{
		Interface:  kubeClient,
		khulnasoftClient: khulnasoftClient,
	}, nil
}

type client struct {
	kubernetes.Interface
	khulnasoftClient khulnasoftclient.Interface
}

func (c *client) Khulnasoft() khulnasoftclient.Interface {
	return c.khulnasoftClient
}
