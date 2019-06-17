package typesnew

type BaseService struct {
	// Docker image of the service
	Image string `json:"image"`
	// Extra arguments that are added to the services
	ExtraArgs map[string]string `json:"extraArgs"`
	// Extra binds added to the nodes
	ExtraBinds []string `json:"extraBinds"`
	// this is to provide extra env variable to the docker container running kubernetes service
	ExtraEnv []string `json:"extraEnv"`
}

type ETCDService struct {
	// Base service properties
	BaseService `json:",inline"`
	// List of etcd urls
	ExternalURLs []string `json:"externalUrls"`
	// External CA certificate
	CACert string `json:"caCert"`
	// External Client certificate
	Cert string `json:"cert"`
	// External Client key
	Key string `json:"key"`
	// External etcd prefix
	Path string `json:"path"`
	// Etcd Recurring snapshot Service
	Snapshot *bool `json:"snapshot"`
	// Etcd snapshot Retention period
	Retention string `json:"retention"`
	// Etcd snapshot Creation period
	Creation string `json:"creation"`
	// Backup backend for etcd snapshots, used by zke only
	BackupConfig *BackupConfig `json:"backupConfig"`
}

type KubeAPIService struct {
	// Base service properties
	BaseService `json:",inline"`
	// Virtual IP range that will be used by Kubernetes services
	ServiceClusterIPRange string `json:"serviceClusterIpRange"`
	// Port range for services defined with NodePort type
	ServiceNodePortRange string `json:"serviceNodePortRange"`
	// Enabled/Disable PodSecurityPolicy
	PodSecurityPolicy bool `json:"podSecurityPolicy"`
	// Enable/Disable AlwaysPullImages admissions plugin
	AlwaysPullImages bool `json:"alwaysPullImages"`
}

type KubeControllerService struct {
	// Base service properties
	BaseService `json:",inline"`
	// CIDR Range for Pods in cluster
	ClusterCIDR string `json:"clusterCidr"`
	// Virtual IP range that will be used by Kubernetes services
	ServiceClusterIPRange string `json:"serviceClusterIpRange"`
}

type KubeletService struct {
	// Base service properties
	BaseService `json:",inline"`
	// Domain of the cluster (default: "cluster.local")
	ClusterDomain string `json:"clusterDomain"`
	// The image whose network/ipc namespaces containers in each pod will use
	InfraContainerImage string `json:"infraContainerImage"`
	// Cluster DNS service ip
	ClusterDNSServer string `json:"clusterDnsServer"`
	// Fail if swap is enabled
	FailSwapOn bool `json:"failSwapOn"`
}

type KubeproxyService struct {
	BaseService `json:",inline"`
}

type SchedulerService struct {
	BaseService `json:",inline"`
}

type BackupConfig struct {
	// Backup interval in hours
	IntervalHours int `json:"intervalHours"`
	// Number of backups to keep
	Retention int `json:"retention"`
}

type AuthnConfig struct {
	// Authentication strategy that will be used in kubernetes cluster
	Strategy string `json:"strategy"`
	// List of additional hostnames and IPs to include in the api server PKI cert
	SANs []string `json:"sans"`
	// Webhook configuration options
	Webhook *AuthWebhookConfig `json:"webhook"`
}

type AuthzConfig struct {
	// Authorization mode used by kubernetes
	Mode string `json:"mode"`
	// Authorization mode options
	Options map[string]string `json:"options"`
}

type AuthWebhookConfig struct {
	// ConfigFile is a multiline string that represent a custom webhook config file
	ConfigFile string `json:"configFile"`
	// CacheTimeout controls how long to cache authentication decisions
	CacheTimeout string `json:"cacheTimeout"`
}
