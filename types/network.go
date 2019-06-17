package types

type NetworkConfig struct {
	Plugin  string        `yaml:"plugin" json:"plugin"`
	Iface   string        `yaml:"iface" json:"iface"`
	DNS     DNSConfig     `yaml:"dns" json:"dns"`
	Ingress IngressConfig `yaml:"ingress" json:"ingress"`
}

type DNSConfig struct {
	// DNS provider
	Provider string `yaml:"provider" json:"provider" norman:"coredns"`
	// Upstream nameservers
	UpstreamNameservers []string `yaml:"upstreamnameservers" json:"upstreamnameservers"`
	// ReverseCIDRs
	ReverseCIDRs []string `yaml:"reversecidrs" json:"reversecidrs"`
	// NodeSelector key pair
	NodeSelector map[string]string `yaml:"node_selector" json:"nodeSelector"`
}

type IngressConfig struct {
	// Ingress controller type used by kubernetes
	Provider string `yaml:"provider" json:"provider"`
	// Ingress controller options
	Options map[string]string `yaml:"options" json:"options"`
	// NodeSelector key pair
	NodeSelector map[string]string `yaml:"node_selector" json:"nodeSelector"`
	// Ingress controller extra arguments
	ExtraArgs map[string]string `yaml:"extra_args" json:"extraArgs"`
}
