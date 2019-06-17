package types

type ZKEConfig struct {
	ClusterName string `yaml:"name" json:"name"`
	// ZKE config global options
	Option ZKEConfigOption `yaml:"option" json:"option"`
	// Kubernetes nodes
	Nodes []ZKEConfigNode `yaml:"nodes" json:"nodes"`
	// Kubernetes core components
	Core ZKEConfigCore `yaml:"core" json:"core"`
	// Network configuration used in the kubernetes cluster (flannel, calico, dns and ingress)
	Network ZKEConfigNetwork `yaml:"network" json:"network"`
	// List images used when ZKE create one kubernetes cluster
	SystemImage ZKESystemImages `yaml:"system_images" json:"systemImages"`
	// List of private registries and their credentials
	PrivateRegistries []PrivateRegistry `yaml:"private_registries" json:"privateRegistries"`
	// ZKE Config version
	Version string `yaml:"version" json:"version"`
}

type ZKEConfigOption struct {
	SSHUser               string `yaml:"ssh_user" json:"sshUser"`
	SSHKey                string `yaml:"ssh_key" json:"sshKey"`
	SSHKeyPath            string `yaml:"ssh_key_path" json:"sshKeyPath"`
	SSHPort               string `yaml:"ssh_port" json:"sshPort"`
	DockerSocket          string `yaml:"docker_socket" json:"dockerSocket"`
	KubernetesVersion     string `yaml:"kubetnetes_version" json:"kubernetesVersion"`
	IgnoreDockerVersion   bool   `yaml:"ignore_docker_version" json:"ignoreDockerVersion"`
	ClusterCidr           string `yaml:"cluster_cidr" json:"clusterCidr"`
	ServiceClusterIpRange string `yaml:"service_cluster_ip_range" json:"serviceClusterIpRange"`
	ClusterDomain         string `yaml:"cluster_domain" json:"clusterDomain"`
	PrefixPath            string `yaml:"prefix_path" json:"prefixPath"`
}

type ZKEConfigCore struct {
	Etcd           ETCDService           `yaml:"etcd" json:"etcd"`
	KubeAPI        KubeAPIService        `yaml:"kube_api" json:"kubeApi"`
	KubeController KubeControllerService `yaml:"kube_controller" json:"kubeController"`
	Scheduler      SchedulerService      `yaml:"scheduler" json:"scheduler"`
	Kubelet        KubeletService        `yaml:"kubelet" json:"kubelet"`
	Kubeproxy      KubeproxyService      `yaml:"kube_proxy" json:"kubeproxy"`
	Authentication AuthnConfig           `yaml:"authentication" json:"authentication"`
	Authorization  AuthzConfig           `yaml:"authorization" json:"authorization"`
}

type ZKEConfigNetwork struct {
	// Network Plugin That will be used in kubernetes cluster
	Plugin string `yaml:"plugin" json:"plugin"`
	// Plugin options to configure network properties
	Iface   string        `yaml:"iface" json:"iface"`
	DNS     DNSConfig     `yaml:"dns" json:"dns"`
	Ingress IngressConfig `yaml:"ingress" json:"ingress"`
}
