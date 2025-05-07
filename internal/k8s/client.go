package k8s

import (
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	gatewayclientset "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

type Client struct {
	Config        *rest.Config
	GatewayClient *gatewayclientset.Clientset
}

func NewClient() (*Client, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	config, err := BuildConfig(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build k8s config: %w", err)
	}

	gatewayClient, err := gatewayclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create gateway clientset from config: %w", err)
	}
	log.Println("Successfully created Kubernetes Gateway API client.")

	return &Client{
		Config:        config,
		GatewayClient: gatewayClient,
	}, nil
}

// BuildConfig creates a Kubernetes client configuration from the given kubeconfig path
// or in-cluster configuration if path is empty.
func BuildConfig(kubeconfigPath string) (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build k8s rest.Config: %w", err)
	}

	log.Printf("Successfully built Kubernetes rest.Config. Kubeconfig: '%s', In-cluster: %t\n",
		kubeconfigPath, kubeconfigPath == "")

	return config, nil
}
