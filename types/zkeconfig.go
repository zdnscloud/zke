package types

type ZKEConfig struct {
	ClusterName string `json:"name"`
	// ZKE config global options
	Option ZKEConfigOption `json:"option"`
	// Kubernetes nodes
	Nodes []ZKEConfigNode `json:"nodes"`
	// Kubernetes core components
	Core ZKEConfigCore `json:"core"`
	// Network configuration used in the kubernetes cluster (flannel, calico, dns and ingress)
	Network ZKEConfigNetwork `json:"network"`
	// List images used when ZKE create one kubernetes cluster
	SystemImage ZKESystemImages `json:"systemImages"`
	// List of private registries and their credentials
	PrivateRegistries []PrivateRegistry `json:"privateRegistries"`
	// ZKE Config version
	Version string `json:"version"`
}

type ZKEConfigOption struct {
	SSHUser               string `json:"sshUser"`
	SSHKey                string `json:"sshKey"`
	SSHPort               string `json:"sshPort"`
	DockerSocket          string `json:"dockerSocket"`
	KubernetesVersion     string `json:"kubernetesVersion"`
	IgnoreDockerVersion   bool   `json:"ignoreDockerVersion"`
	ClusterCidr           string `json:"clusterCidr"`
	ServiceClusterIpRange string `json:"serviceClusterIpRange"`
	ClusterDomain         string `json:"clusterDomain"`
}

type ZKEConfigCore struct {
	Etcd           ETCDService           `json:"etcd"`
	KubeAPI        KubeAPIService        `json:"kubeApi"`
	KubeController KubeControllerService `json:"kubeController"`
	Scheduler      SchedulerService      `json:"scheduler"`
	Kubelet        KubeletService        `json:"kubelet"`
	Kubeproxy      KubeproxyService      `json:"kubeproxy"`
	Authentication AuthnConfig           `json:"authentication"`
	Authorization  AuthzConfig           `json:"authorization"`
}

type ZKEConfigNetwork struct {
	// Network Plugin That will be used in kubernetes cluster
	Plugin string `json:"plugin"`
	// Plugin options to configure network properties
	Iface   string        `json:"iface"`
	DNS     DNSConfig     `json:"dns"`
	Ingress IngressConfig `json:"ingress"`
}
