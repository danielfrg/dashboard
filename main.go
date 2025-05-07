// dashboard/main.go
package main

import (
	"dashboard/internal/api"
	"dashboard/internal/k8s"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Application struct {
	k8sClient *k8s.Client
}

// Simple CORS middleware to allow requests from any origin
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin in development
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
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

	// Create file server handler for the web directory
	webDir := filepath.Join(".", "frontend/dist")
	fileServer := http.FileServer(http.Dir(webDir))

	// Create a custom mux to apply our middleware
	mux := http.NewServeMux()

	// API routes with CORS middleware
	mux.HandleFunc("/api/routes", func(w http.ResponseWriter, r *http.Request) {
		httpRouteHandler.ServeRoutes(w, r, "default")
	})

	// Serve static files from the web directory
	mux.Handle("/assets/", fileServer)
	mux.Handle("/images/", fileServer)
	mux.Handle("/favicon.ico", fileServer)

	// Root handler that serves index.html for root path but falls back to API handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Try to serve index.html from web directory
			indexPath := filepath.Join(webDir, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				http.ServeFile(w, r, indexPath)
				return
			}
			// Fall back to API index handler if index.html doesn't exist
			handleIndex(w, r)
			return
		}

		// Try to serve the file from the web directory
		fileServer.ServeHTTP(w, r)
	})

	// Apply CORS middleware to all requests
	handler := corsMiddleware(mux)

	port := "8080"
	host := "localhost"
	fmt.Printf("Starting HTTP server on port %s...\n", port)
	fmt.Printf("Serving static files from: %s\n", webDir)
	log.Fatal(http.ListenAndServe(host + ":" + port, handler))
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
