package typesnew

type PrivateRegistry struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
	CAcert   string `json:"caCert"`
}

type ZKEConfigNode struct {
	NodeName string `json:"name"`
	Address  string `json:"address"`
	// Optional - Internal address that will be used for components communication
	InternalAddress string `json:"internalAddress"`
	// Node role in kubernetes cluster (controlplane, worker, etcd, storage or edge)
	Roles []string `json:"roles"`
	// Optional - Hostname of the node
	HostnameOverride string `json:"hostnameOverride"`
	// SSH config
	User         string            `json:"sshUser"`
	Port         string            `json:"sshPort"`
	SSHKey       string            `json:"sshKey"`
	DockerSocket string            `json:"dockerSocket"`
	Labels       map[string]string `json:"labels"`
}
