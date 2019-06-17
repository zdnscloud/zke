package typesnew

type DNSConfig struct {
	// DNS provider
	Provider string `json:"provider"`
	// Upstream nameservers
	UpstreamNameservers []string `json:"upstreamnameservers"`
	// ReverseCIDRs
	ReverseCIDRs []string `json:"reversecidrs"`
	// NodeSelector key pair
	NodeSelector map[string]string `json:"nodeSelector"`
}

type IngressConfig struct {
	// Ingress controller type used by kubernetes
	Provider string `json:"provider"`
	// Ingress controller options
	Options map[string]string `json:"options"`
	// NodeSelector key pair
	NodeSelector map[string]string `json:"nodeSelector"`
	// Ingress controller extra arguments
	ExtraArgs map[string]string `json:"extraArgs"`
}
