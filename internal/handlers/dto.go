package handlers

type ParentRef struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
}

type HTTPRouteInfo struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Hostnames   []string          `json:"hostnames,omitempty"`
	ParentRefs  []ParentRef       `json:"parentRefs,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type RoutesResponse struct {
	QueryNamespace string          `json:"queryNamespace"`
	Count          int             `json:"count"`
	Routes         []HTTPRouteInfo `json:"routes"`
}
