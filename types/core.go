package types

type BaseService struct {
	// Docker image of the service
	Image string `yaml:"image" json:"image,omitempty"`
	// Extra arguments that are added to the services
	ExtraArgs map[string]string `yaml:"extra_args" json:"extraArgs,omitempty"`
	// Extra binds added to the nodes
	ExtraBinds []string `yaml:"extra_binds" json:"extraBinds,omitempty"`
	// this is to provide extra env variable to the docker container running kubernetes service
	ExtraEnv []string `yaml:"extra_env" json:"extraEnv,omitempty"`
}

type ETCDService struct {
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
	// List of etcd urls
	ExternalURLs []string `yaml:"external_urls" json:"externalUrls,omitempty"`
	// External CA certificate
	CACert string `yaml:"ca_cert" json:"caCert,omitempty"`
	// External Client certificate
	Cert string `yaml:"cert" json:"cert,omitempty"`
	// External Client key
	Key string `yaml:"key" json:"key,omitempty"`
	// External etcd prefix
	Path string `yaml:"path" json:"path,omitempty"`
	// Etcd Recurring snapshot Service
	Snapshot *bool `yaml:"snapshot" json:"snapshot,omitempty" norman:"default=true"`
	// Etcd snapshot Retention period
	Retention string `yaml:"retention" json:"retention,omitempty" norman:"default=72h"`
	// Etcd snapshot Creation period
	Creation string `yaml:"creation" json:"creation,omitempty" norman:"default=12h"`
	// Backup backend for etcd snapshots, used by zke only
	BackupConfig *BackupConfig `yaml:"backup_config" json:"backupConfig,omitempty"`
}

type KubeAPIService struct {
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
	// Virtual IP range that will be used by Kubernetes services
	ServiceClusterIPRange string `yaml:"service_cluster_ip_range" json:"serviceClusterIpRange,omitempty"`
	// Port range for services defined with NodePort type
	ServiceNodePortRange string `yaml:"service_node_port_range" json:"serviceNodePortRange,omitempty" norman:"default=30000-32767"`
	// Enabled/Disable PodSecurityPolicy
	PodSecurityPolicy bool `yaml:"pod_security_policy" json:"podSecurityPolicy,omitempty"`
	// Enable/Disable AlwaysPullImages admissions plugin
	AlwaysPullImages bool `yaml:"always_pull_images" json:"always_pull_images,omitempty"`
}

type KubeControllerService struct {
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
	// CIDR Range for Pods in cluster
	ClusterCIDR string `yaml:"cluster_cidr" json:"clusterCidr,omitempty"`
	// Virtual IP range that will be used by Kubernetes services
	ServiceClusterIPRange string `yaml:"service_cluster_ip_range" json:"serviceClusterIpRange,omitempty"`
}

type KubeletService struct {
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
	// Domain of the cluster (default: "cluster.local")
	ClusterDomain string `yaml:"cluster_domain" json:"clusterDomain,omitempty"`
	// The image whose network/ipc namespaces containers in each pod will use
	InfraContainerImage string `yaml:"infra_container_image" json:"infraContainerImage,omitempty"`
	// Cluster DNS service ip
	ClusterDNSServer string `yaml:"cluster_dns_server" json:"clusterDnsServer,omitempty"`
	// Fail if swap is enabled
	FailSwapOn bool `yaml:"fail_swap_on" json:"failSwapOn,omitempty"`
}

type KubeproxyService struct {
	BaseService `yaml:",inline" json:",inline"`
}

type SchedulerService struct {
	BaseService `yaml:",inline" json:",inline"`
}

type BackupConfig struct {
	// Backup interval in hours
	IntervalHours int `yaml:"interval_hours" json:"intervalHours,omitempty" norman:"default=12"`
	// Number of backups to keep
	Retention int `yaml:"retention" json:"retention,omitempty" norman:"default=6"`
}

type AuthnConfig struct {
	// Authentication strategy that will be used in kubernetes cluster
	Strategy string `yaml:"strategy" json:"strategy,omitempty" norman:"default=x509"`
	// List of additional hostnames and IPs to include in the api server PKI cert
	SANs []string `yaml:"sans" json:"sans,omitempty"`
	// Webhook configuration options
	Webhook *AuthWebhookConfig `yaml:"webhook" json:"webhook,omitempty"`
}

type AuthzConfig struct {
	// Authorization mode used by kubernetes
	Mode string `yaml:"mode" json:"mode,omitempty"`
	// Authorization mode options
	Options map[string]string `yaml:"options" json:"options,omitempty"`
}

type AuthWebhookConfig struct {
	// ConfigFile is a multiline string that represent a custom webhook config file
	ConfigFile string `yaml:"config_file" json:"configFile,omitempty"`
	// CacheTimeout controls how long to cache authentication decisions
	CacheTimeout string `yaml:"cache_timeout" json:"cacheTimeout,omitempty"`
}
