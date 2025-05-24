package kube

import (
	"fmt"

	loftclient "dev.khulnasoft.com/api/pkg/clientset/versioned"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Interface interface {
	kubernetes.Interface
	Loft() loftclient.Interface
}

func NewForConfig(c *rest.Config) (Interface, error) {
	kubeClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, fmt.Errorf("create kube client: %w", err)
	}

	loftClient, err := loftclient.NewForConfig(c)
	if err != nil {
		return nil, fmt.Errorf("create loft client: %w", err)
	}

	return &client{
		Interface:  kubeClient,
		loftClient: loftClient,
	}, nil
}

type client struct {
	kubernetes.Interface
	loftClient loftclient.Interface
}

func (c *client) Loft() loftclient.Interface {
	return c.loftClient
}
