// Package k8sclient drives Kubernetes (via client-go) to create and tear
// down the Postgres and app pods/services used for each benchmark run,
// mirroring the previous src/k8s/*.ts (which used @kubernetes/client-node).
//
// Because this talks to whatever context the kubeconfig points at, the exact
// same code runs against a local minikube cluster today and a real cloud
// cluster (EKS/GKE/AKS) later — only the kubeconfig changes.
package k8sclient

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Clientset *kubernetes.Clientset
	Namespace string
}

// New builds a Client from the default kubeconfig loading rules (respects
// $KUBECONFIG and ~/.kube/config), the same default lookup client-node's
// loadFromDefault() used.
func New(namespace string) (*Client, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("loading kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	return &Client{Clientset: clientset, Namespace: namespace}, nil
}
