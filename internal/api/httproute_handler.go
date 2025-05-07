package api

import (
	"dashboard/internal/k8s"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func (h *HTTPRouteHandler) ServeRoutes(w http.ResponseWriter, r *http.Request, queryNamespace string) {
	rawRoutes, err := k8s.ListHTTPRoutes(h.gatewayClient, queryNamespace)
	if err != nil {
		log.Printf("Error listing raw HTTPRoutes for queryNamespace '%s': %v", queryNamespace, err)
		http.Error(w, fmt.Sprintf("Error listing HTTPRoutes: %v", err), http.StatusInternalServerError)
		return
	}

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
			Name:      route.Name,
			Namespace: route.Namespace,
		}

		if len(route.Spec.Hostnames) > 0 {
			routeInfo.Hostnames = make([]string, len(route.Spec.Hostnames))
			for i,hostname := range route.Spec.Hostnames {
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
