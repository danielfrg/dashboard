package k8s

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayclientset "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

func ListHTTPRoutes(client *gatewayclientset.Clientset, namespace string) ([]gatewayv1.HTTPRoute, error) {
	if client == nil {
		return nil, fmt.Errorf("kubernetes gateway client is nil")
	}

	effectiveNamespace := namespace
	if effectiveNamespace == "" {
		effectiveNamespace = metav1.NamespaceAll
	}

	log.Printf("Listing HTTPRoutes in namespace: %q using provided client\n", effectiveNamespace)

	httpRoutes, err := client.GatewayV1().HTTPRoutes(effectiveNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list HTTPRoutes in namespace %s: %w", effectiveNamespace, err)
	}

	return httpRoutes.Items, nil
}
