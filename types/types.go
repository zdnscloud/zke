package types

type ZcloudKubernetesEngineConfig struct {
	Option ZKEConfigOption `yaml:"option" json:"option"`
	Nodes  []ZKEConfigNode `yaml:"nodes" json:"nodes"`
	// Kubernetes components
	Services          ZKEConfigServices `yaml:"services" json:"services"`
	Network           NetworkConfig     `yaml:"network" json:"network"`
	Authentication    AuthnConfig       `yaml:"authentication" json:"authentication"`
	Authorization     AuthzConfig       `yaml:"authorization" json:"authorization"`
	SystemImages      ZKESystemImages   `yaml:"system_images" json:"systemImages"`
	PrivateRegistries []PrivateRegistry `yaml:"private_registries" json:"privateRegistries"`
	ClusterName       string            `yaml:"cluster_name" json:"name"`
	Monitor           MonitorConfig     `yaml:"monitor" json:"monitor"`
	Version           string            `yaml:"version" json:"version"`
}

type PrivateRegistry struct {
	URL      string `yaml:"url" json:"url"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	CAcert   string `yaml:"ca_cert" json:"caCert"`
}

type ZKEConfigNode struct {
	NodeName string `yaml:"node_name" json:"nodeName"`
	Address  string `yaml:"address" json:"address"`
	// Optional - Internal address that will be used for components communication
	InternalAddress string `yaml:"internal_address" json:"internalAddress"`
	// Node role in kubernetes cluster (controlplane, worker, etcd, storage or edge)
	Role []string `yaml:"role" json:"role"`
	// Optional - Hostname of the node
	HostnameOverride string `yaml:"hostname_override" json:"hostnameOverride"`
	// SSH config
	User         string            `yaml:"user" json:"sshUser"`
	Port         string            `yaml:"port" json:"sshPort"`
	SSHKey       string            `yaml:"ssh_key" json:"sshKey"`
	SSHKeyPath   string            `yaml:"ssh_key_path" json:"sshKeyPath"`
	DockerSocket string            `yaml:"docker_socket" json:"dockerSocket"`
	Labels       map[string]string `yaml:"labels" json:"labels"`
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
