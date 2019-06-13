package types

type ZKEConfig struct {
	// ZKE config global options
	Options ZKEConfigOptions `yaml:"options" json:"options"`
	// Kubernetes nodes
	Nodes []ZKEConfigNode `yaml:"nodes" json:"nodes,omitempty"`
	// Kubernetes core components
	Core ZKEConfigCore `yaml:"core" json:"core,omitempty"`
	// Network configuration used in the kubernetes cluster (flannel, calico, dns and ingress)
	Network ZKEConfigNetwork `yaml:"network" json:"network,omitempty"`
	// Zcloud Resource config (monitor and harbor registry)
	Zcloud ZKEConfigZcloud `yaml:"zcloud" json:"zcloud,omitempty"`
	// List images used when ZKE create one kubernetes cluster
	SystemImages ZKESystemImages `yaml:"system_images" json:"systemImages,omitempty"`
	// List of private registries and their credentials
	PrivateRegistries []PrivateRegistry `yaml:"private_registries" json:"privateRegistries,omitempty"`
	// ZKE Config version
	Version string `yaml:"version" json:"Version"`
}

type ZKEConfigOptions struct {
	SSHKeyPath          string `yaml:"ssh_key_path" json:"sshKeyPath,omitempty"`
	SSHKey              string `yaml:"ssh_key" json:"sshKey,omitempty"`
	SSHPort             string `yaml:"port" json:"sshPort,omitempty"`
	SSHUser             string `yaml:"user" json:"sshUser,omitempty"`
	DockerSocket        string `yaml:"docker_socket" json:"dockerSocket,omitempty"`
	K8sVersion          string `yaml:"kubernetes_version" json:"kubernetesVersion,omitempty"`
	ClusterName         string `yaml:"cluster_name" json:"clusterName,omitempty"`
	PrefixPath          string `yaml:"prefix_path" json:"prefixPath,omitempty"`
	IgnoreDockerVersion bool   `yaml:"ignore_docker_version" json:"ignoreDockerVersion" norman:"default=true"`
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
	Options                map[string]string       `yaml:"options" json:"options,omitempty"`
	CalicoNetworkProvider  *CalicoNetworkProvider  `yaml:",omitempty" json:"calicoNetworkProvider,omitempty"`
	FlannelNetworkProvider *FlannelNetworkProvider `yaml:",omitempty" json:"flannelNetworkProvider,omitempty"`
	DNS                    DNSConfig               `yaml:"dns" json:"dns,omitempty"`
	Ingress                IngressConfig           `yaml:"ingress" json:"ingress,omitempty"`
}

type ZKEConfigZcloud struct {
	Storage  StorageConfig  `yaml:"storage" json:"storage,omitempty"`
	Monitor  MonitorConfig  `yaml:"monitoring" json:"monitoring,omitempty"`
	Registry RegistryConfig `yaml:"registry" json:"registry,omitempty"`
}
