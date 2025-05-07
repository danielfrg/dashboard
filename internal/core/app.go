package core

import (
    "fmt"

    "dashboard/internal/k8s"

    "github.com/pocketbase/pocketbase"
)

// App holds application resources.
type App struct {
    Pb        *pocketbase.PocketBase
    K8sClient *k8s.Client
}

func New() (*App, error) {
    k8sClient, err := k8s.NewClient()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize kubernetes client: %w", err)
    }

    pbApp := pocketbase.New()

    return &App{
        Pb:        pbApp,
        K8sClient: k8sClient,
    }, nil
}
