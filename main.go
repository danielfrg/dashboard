// dashboard/main.go
package main

import (
	"dashboard/internal/api"
	"dashboard/internal/k8s"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Application struct {
	k8sClient *k8s.Client
}

func NewApplication() (*Application, error) {
	k8sClient, err := k8s.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize kubernetes client: %w", err)
	}

	return &Application{
		k8sClient: k8sClient,
	}, nil
}

func main() {
	fmt.Printf("Attempting to use KUBECONFIG env var: '%s', or in-cluster config if not set\n", os.Getenv("KUBECONFIG"))
	app, err := NewApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	httpRouteHandler := api.NewHTTPRouteHandler(app.k8sClient.GatewayClient)

	http.HandleFunc("/", handleIndex)

	http.HandleFunc("/api/routes", func(w http.ResponseWriter, r *http.Request) {
		httpRouteHandler.ServeRoutes(w, r, "default")
	})

	port := "8080"
	host := "localhost"
	fmt.Printf("Starting HTTP server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(host + ":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Kubernetes Gateway API Dashboard")
	fmt.Fprintln(w, "-------------------------")
	fmt.Fprintln(w, "View HTTPRoutes in your cluster:")
	fmt.Fprintln(w, "/routes.json - View in default namespace (JSON format)")
	fmt.Fprintln(w, "/routes/all.json - View in all namespaces (JSON format)")
}
