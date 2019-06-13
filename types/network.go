package types

type NetworkConfig struct {
	// Network Plugin That will be used in kubernetes cluster
	Plugin string `yaml:"plugin" json:"plugin,omitempty" norman:"default=canal"`
	// Plugin options to configure network properties
	Options map[string]string `yaml:"options" json:"options,omitempty"`
	// CalicoNetworkProvider
	CalicoNetworkProvider *CalicoNetworkProvider `yaml:",omitempty" json:"calicoNetworkProvider,omitempty"`
	// FlannelNetworkProvider
	FlannelNetworkProvider *FlannelNetworkProvider `yaml:",omitempty" json:"flannelNetworkProvider,omitempty"`
}

type CalicoNetworkProvider struct {
	CloudProvider string `json:"cloudProvider"`
}

type FlannelNetworkProvider struct {
	// Alternate cloud interface for flannel
	Iface string `json:"iface"`
}

type DNSConfig struct {
	// DNS provider
	Provider string `yaml:"provider" json:"provider,omitempty" norman:"coredns"`
	// Upstream nameservers
	UpstreamNameservers []string `yaml:"upstreamnameservers" json:"upstreamnameservers,omitempty"`
	// ReverseCIDRs
	ReverseCIDRs []string `yaml:"reversecidrs" json:"reversecidrs,omitempty"`
	// NodeSelector key pair
	NodeSelector map[string]string `yaml:"node_selector" json:"nodeSelector,omitempty"`
}

type IngressConfig struct {
	// Ingress controller type used by kubernetes
	Provider string `yaml:"provider" json:"provider,omitempty" norman:"default=nginx"`
	// Ingress controller options
	Options map[string]string `yaml:"options" json:"options,omitempty"`
	// NodeSelector key pair
	NodeSelector map[string]string `yaml:"node_selector" json:"nodeSelector,omitempty"`
	// Ingress controller extra arguments
	ExtraArgs map[string]string `yaml:"extra_args" json:"extraArgs,omitempty"`
}
