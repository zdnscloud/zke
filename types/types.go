package types

type ZcloudKubernetesEngineConfig struct {
	// Kubernetes nodes
	Nodes []ZKEConfigNode `yaml:"nodes" json:"nodes,omitempty"`
	// Kubernetes components
	Services ZKEConfigServices `yaml:"services" json:"services,omitempty"`
	Network  NetworkConfig     `yaml:"network" json:"network,omitempty"`
	Storage  StorageConfig     `yaml:"storage" json:"storage,omitempty"`
	// Authentication configuration used in the cluster (default: x509)
	Authentication AuthnConfig `yaml:"authentication" json:"authentication,omitempty"`
	// Authorization mode configuration used in the cluster
	Authorization AuthzConfig `yaml:"authorization" json:"authorization,omitempty"`
	// List of images used internally for proxy, cert downlaod and kubedns
	SystemImages ZKESystemImages `yaml:"system_images" json:"systemImages,omitempty"`
	// SSH Private Key Path
	SSHKeyPath   string `yaml:"ssh_key_path" json:"sshKeyPath,omitempty"`
	SSHKey       string `yaml:"ssh_key" json:"sshKey,omitempty"`
	SSHPort      string `yaml:"port" json:"sshPort,omitempty"`
	SSHUser      string `yaml:"user" json:"sshUser,omitempty"`
	DockerSocket string `yaml:"docker_socket" json:"dockerSocket,omitempty"`
	// Enable/disable strict docker version checking
	IgnoreDockerVersion bool `yaml:"ignore_docker_version" json:"ignoreDockerVersion" norman:"default=true"`
	// Kubernetes version to use (if kubernetes image is specifed, image version takes precedence)
	Version string `yaml:"kubernetes_version" json:"kubernetesVersion,omitempty"`
	// List of private registries and their credentials
	PrivateRegistries []PrivateRegistry `yaml:"private_registries" json:"privateRegistries,omitempty"`
	Ingress           IngressConfig     `yaml:"ingress" json:"ingress,omitempty"`
	ClusterName       string            `yaml:"cluster_name" json:"clusterName,omitempty"`
	PrefixPath        string            `yaml:"prefix_path" json:"prefixPath,omitempty"`
	Monitor           MonitorConfig     `yaml:"monitoring" json:"monitoring,omitempty"`
	DNS               DNSConfig         `yaml:"dns" json:"dns,omitempty"`
	// Harbor Registry Config
	Registry RegistryConfig `yaml:"registry" json:"registry,omitempty"`
	// ZKE config version
	ConfigVersion string `yaml:"config_version" json:"config_version"`
}

type PrivateRegistry struct {
	URL      string `yaml:"url" json:"url,omitempty"`
	User     string `yaml:"user" json:"user,omitempty"`
	Password string `yaml:"password" json:"password,omitempty" norman:"type=password"`
	CAcert   string `yaml:"ca_cert" json:"cacert,omitempty"`
}

type ZKEConfigNode struct {
	NodeName string `yaml:"nodeName,omitempty" json:"nodeName,omitempty" norman:"type=reference[node]"`
	Address  string `yaml:"address" json:"address,omitempty"`
	// Optional - Internal address that will be used for components communication
	InternalAddress string `yaml:"internal_address" json:"internalAddress,omitempty"`
	// Node role in kubernetes cluster (controlplane, worker, etcd, storage or edge)
	Role []string `yaml:"role" json:"role,omitempty" norman:"type=array[enum],options=etcd|worker|controlplane"`
	// Optional - Hostname of the node
	HostnameOverride string `yaml:"hostname_override" json:"hostnameOverride,omitempty"`
	// SSH config
	User         string            `yaml:"user" json:"user,omitempty"`
	Port         string            `yaml:"port" json:"port,omitempty"`
	SSHKey       string            `yaml:"ssh_key" json:"sshKey,omitempty" norman:"type=password"`
	SSHKeyPath   string            `yaml:"ssh_key_path" json:"sshKeyPath,omitempty"`
	DockerSocket string            `yaml:"docker_socket" json:"dockerSocket,omitempty"`
	Labels       map[string]string `yaml:"labels" json:"labels,omitempty"`
}

type ZKEConfigNodePlan struct {
	Address string `json:"address,omitempty"`
	// map of named processes that should run on the node
	Processes   map[string]Process `json:"processes,omitempty"`
	PortChecks  []PortCheck        `json:"portChecks,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
	Labels      map[string]string  `json:"labels,omitempty"`
}

type ZKEConfigServices struct {
	Etcd           ETCDService           `yaml:"etcd" json:"etcd,omitempty"`
	KubeAPI        KubeAPIService        `yaml:"kube-api" json:"kubeApi,omitempty"`
	KubeController KubeControllerService `yaml:"kube-controller" json:"kubeController,omitempty"`
	Scheduler      SchedulerService      `yaml:"scheduler" json:"scheduler,omitempty"`
	Kubelet        KubeletService        `yaml:"kubelet" json:"kubelet,omitempty"`
	Kubeproxy      KubeproxyService      `yaml:"kubeproxy" json:"kubeproxy,omitempty"`
}

type Process struct {
	Name    string   `json:"name,omitempty"`
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
	Env     []string `json:"env,omitempty"`
	Image   string   `json:"image,omitempty"`
	//AuthConfig for image private registry
	ImageRegistryAuthConfig string `json:"imageRegistryAuthConfig,omitempty"`
	// Process docker image VolumesFrom
	VolumesFrom []string `json:"volumesFrom,omitempty"`
	// Process docker container bind mounts
	Binds         []string    `json:"binds,omitempty"`
	NetworkMode   string      `json:"networkMode,omitempty"`
	RestartPolicy string      `json:"restartPolicy,omitempty"`
	PidMode       string      `json:"pidMode,omitempty"`
	Privileged    bool        `json:"privileged,omitempty"`
	HealthCheck   HealthCheck `json:"healthCheck,omitempty"`
	// Process docker container Labels
	Labels map[string]string `json:"labels,omitempty"`
	// Process docker publish container's port to host
	Publish []string `json:"publish,omitempty"`
}

type HealthCheck struct {
	URL string `json:"url,omitempty"`
}

type PortCheck struct {
	Address  string `json:"address,omitempty"`
	Port     int    `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

type KubernetesServicesOptions struct {
	KubeAPI        map[string]string `json:"kubeapi"`
	Kubelet        map[string]string `json:"kubelet"`
	Kubeproxy      map[string]string `json:"kubeproxy"`
	KubeController map[string]string `json:"kubeController"`
	Scheduler      map[string]string `json:"scheduler"`
}
