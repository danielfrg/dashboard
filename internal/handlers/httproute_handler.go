package handlers

import (
	"dashboard/internal/k8s" // Corrected import path
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4" // Use Echo context
	gatewayclientset "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

type HTTPRouteHandler struct {
	gatewayClient *gatewayclientset.Clientset
}

func NewHTTPRouteHandler(gatewayClient *gatewayclientset.Clientset) *HTTPRouteHandler {
	return &HTTPRouteHandler{
		gatewayClient: gatewayClient,
	}
}

// ServeRoutes now accepts echo.Context
func (h *HTTPRouteHandler) ServeRoutes(c echo.Context) error {
	// Example: Get namespace from query param if needed
	// queryNamespace := c.QueryParam("namespace")
	// For now, hardcode "default" as originally intended
	queryNamespace := "default"

	rawRoutes, err := k8s.ListHTTPRoutes(h.gatewayClient, queryNamespace)
	if err != nil {
		log.Printf("Error listing raw HTTPRoutes for queryNamespace '%s': %v", queryNamespace, err)
		// Use Echo's error handling
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Error listing HTTPRoutes: %v", err)})
	}

	// --- Logic to build the response (same as before) ---
	displayQueryNamespace := queryNamespace
	if displayQueryNamespace == "" {
		displayQueryNamespace = "all"
	}

	response := RoutesResponse{
		QueryNamespace: fmt.Sprintf("%s namespace(s)", displayQueryNamespace),
		Count:          len(rawRoutes),
		Routes:         make([]HTTPRouteInfo, 0, len(rawRoutes)),
	}

	for _, route := range rawRoutes {
		routeInfo := HTTPRouteInfo{
			Name:        route.Name,
			Namespace:   route.Namespace,
			Annotations: route.Annotations,
		}

		if len(route.Spec.Hostnames) > 0 {
			routeInfo.Hostnames = make([]string, len(route.Spec.Hostnames))
			for i, hostname := range route.Spec.Hostnames {
				routeInfo.Hostnames[i] = string(hostname)
			}
		}

		if len(route.Spec.ParentRefs) > 0 {
			routeInfo.ParentRefs = make([]ParentRef, len(route.Spec.ParentRefs))
			for i, pRef := range route.Spec.ParentRefs {
				parentNs := route.Namespace
				if pRef.Namespace != nil && string(*pRef.Namespace) != "" {
					parentNs = string(*pRef.Namespace)
				}
				kind := "Unknown"
				if pRef.Kind != nil {
					kind = string(*pRef.Kind)
				}
				routeInfo.ParentRefs[i] = ParentRef{
					Name:      string(pRef.Name),
					Kind:      kind,
					Namespace: parentNs,
				}
			}
		}
		response.Routes = append(response.Routes, routeInfo)
	}
	// --- End of response building logic ---

	// Use Echo's JSON response method
	return c.JSON(http.StatusOK, response)
}
