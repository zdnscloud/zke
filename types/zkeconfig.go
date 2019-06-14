package types

type ZKEConfig struct {
	ClusterName string `yaml:"cluster_name" json:"name,omitempty"`
	// ZKE config global options
	Option ZKEConfigOption `yaml:"options" json:"options"`
	// Kubernetes nodes
	Nodes []ZKEConfigNode `yaml:"nodes" json:"nodes,omitempty"`
	// Kubernetes core components
	Core ZKEConfigCore `yaml:"core" json:"core,omitempty"`
	// Network configuration used in the kubernetes cluster (flannel, calico, dns and ingress)
	Network ZKEConfigNetwork `yaml:"network" json:"network,omitempty"`
	// List images used when ZKE create one kubernetes cluster
	SystemImage ZKESystemImages `yaml:"system_images" json:"systemImages,omitempty"`
	// List of private registries and their credentials
	PrivateRegistries []PrivateRegistry `yaml:"private_registries" json:"privateRegistries,omitempty"`
	// ZKE Config version
	Version string `yaml:"version" json:"Version"`
}

type ZKEConfigOption struct {
	SSHUser               string `yaml:"user" json:"sshUser,omitempty"`
	SSHKey                string `yaml:"ssh_key" json:"sshKey,omitempty"`
	SSHPort               string `yaml:"port" json:"sshPort,omitempty"`
	DockerSocket          string `yaml:"docker_socket" json:"dockerSocket,omitempty"`
	KubernetesVersion     string `yaml:"kubernetes_version" json:"kubernetesVersion,omitempty"`
	IgnoreDockerVersion   bool   `yaml:"ignore_docker_version" json:"ignoreDockerVersion" norman:"default=true"`
	ClusterCidr           string `yaml:"cluster_cidr" json:"clusterCidr"`
	ServiceClusterIpRange string `yaml:"service_cluster_iprange" json:"serviceClusterIpRange"`
	ClusterDomain         string `yaml:"cluster_domain" json:"clusterDomain"`
}

type ZKEConfigCore struct {
	Etcd           ETCDService           `yaml:"etcd" json:"etcd,omitempty"`
	KubeAPI        KubeAPIService        `yaml:"kube-api" json:"kubeApi,omitempty"`
	KubeController KubeControllerService `yaml:"kube-controller" json:"kubeController,omitempty"`
	Scheduler      SchedulerService      `yaml:"scheduler" json:"scheduler,omitempty"`
	Kubelet        KubeletService        `yaml:"kubelet" json:"kubelet,omitempty"`
	Kubeproxy      KubeproxyService      `yaml:"kubeproxy" json:"kubeproxy,omitempty"`
	Authentication AuthnConfig           `yaml:"authentication" json:"authentication,omitempty"`
	Authorization  AuthzConfig           `yaml:"authorization" json:"authorization,omitempty"`
}

type ZKEConfigNetwork struct {
	// Network Plugin That will be used in kubernetes cluster
	Plugin string `yaml:"plugin" json:"plugin,omitempty" norman:"default=canal"`
	// Plugin options to configure network properties
	Iface   string        `yaml:"iface" json:"iface"`
	DNS     DNSConfig     `yaml:"dns" json:"dns,omitempty"`
	Ingress IngressConfig `yaml:"ingress" json:"ingress,omitempty"`
}
