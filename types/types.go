package types

type ZcloudKubernetesEngineConfig struct {
	// Kubernetes nodes
	Nodes []ZKEConfigNode `yaml:"nodes" json:"nodes,omitempty"`
	// Kubernetes components
	Services ZKEConfigServices `yaml:"services" json:"services,omitempty"`
	// Network configuration used in the kubernetes cluster (flannel, calico)
	Network NetworkConfig `yaml:"network" json:"network,omitempty"`
	// Storage configuration used in the kubernetes cluster (lvm)
	Storage StorageConfig `yaml:"storage" json:"storage,omitempty"`
	// Authentication configuration used in the cluster (default: x509)
	Authentication AuthnConfig `yaml:"authentication" json:"authentication,omitempty"`
	// YAML manifest for user provided addons to be deployed on the cluster
	Addons string `yaml:"addons" json:"addons,omitempty"`
	// List of urls or paths for addons
	AddonsInclude []string `yaml:"addons_include" json:"addonsInclude,omitempty"`
	// List of images used internally for proxy, cert downlaod and kubedns
	SystemImages ZKESystemImages `yaml:"system_images" json:"systemImages,omitempty"`
	// SSH Private Key Path
	SSHKeyPath   string `yaml:"ssh_key_path" json:"sshKeyPath,omitempty"`
	SSHPort      string `yaml:"port" json:"sshPort,omitempty"`
	SSHUser      string `yaml:"user" json:"sshUser,omitempty"`
	DockerSocket string `yaml:"docker_socket" json:"dockerSocket,omitempty"`
	// SSH Certificate Path
	SSHCertPath string `yaml:"ssh_cert_path" json:"sshCertPath,omitempty"`
	// SSH Agent Auth enable
	SSHAgentAuth bool `yaml:"ssh_agent_auth" json:"sshAgentAuth"`
	// Authorization mode configuration used in the cluster
	Authorization AuthzConfig `yaml:"authorization" json:"authorization,omitempty"`
	// Enable/disable strict docker version checking
	IgnoreDockerVersion bool `yaml:"ignore_docker_version" json:"ignoreDockerVersion" norman:"default=true"`
	// Kubernetes version to use (if kubernetes image is specifed, image version takes precedence)
	Version string `yaml:"kubernetes_version" json:"kubernetesVersion,omitempty"`
	// List of private registries and their credentials
	PrivateRegistries []PrivateRegistry `yaml:"private_registries" json:"privateRegistries,omitempty"`
	// Ingress controller used in the cluster
	Ingress IngressConfig `yaml:"ingress" json:"ingress,omitempty"`
	// Cluster Name used in the kube config
	ClusterName string `yaml:"cluster_name" json:"clusterName,omitempty"`
	// Cloud Provider options
	CloudProvider CloudProvider `yaml:"cloud_provider" json:"cloudProvider,omitempty"`
	// kubernetes directory path
	PrefixPath string `yaml:"prefix_path" json:"prefixPath,omitempty"`
	// Timeout in seconds for status check on addon deployment jobs
	AddonJobTimeout int `yaml:"addon_job_timeout" json:"addonJobTimeout,omitempty" norman:"default=30"`
	// Monitoring Config
	Monitor MonitorConfig `yaml:"monitoring" json:"monitoring,omitempty"`
	// RestoreCluster flag
	Restore RestoreConfig `yaml:"restore" json:"restore,omitempty"`
	// DNS Config
	DNS DNSConfig `yaml:"dns" json:"dns,omitempty"`
	// Harbor Registry Config
	Registry      RegistryConfig `yaml:"registry" json:"registry,omitempty"`
	ConfigVersion string         `yaml:"config_version" json:"config_version"`
}

type PrivateRegistry struct {
	// URL for the registry
	URL string `yaml:"url" json:"url,omitempty"`
	// User name for registry acces
	User string `yaml:"user" json:"user,omitempty"`
	// Password for registry access
	Password string `yaml:"password" json:"password,omitempty" norman:"type=password"`
	// Default registry
	// CAcert string `yaml:"ca_cert" json:"ca_cert",omitempty`
}

type ZKESystemImages struct {
	// etcd image
	Etcd string `yaml:"etcd" json:"etcd,omitempty"`
	// Alpine image
	Alpine string `yaml:"alpine" json:"alpine,omitempty"`
	// zke-nginx-proxy image
	NginxProxy string `yaml:"nginx_proxy" json:"nginxProxy,omitempty"`
	// zke-cert-deployer image
	CertDownloader string `yaml:"cert_downloader" json:"certDownloader,omitempty"`
	// zke-service-sidekick image
	KubernetesServicesSidecar string `yaml:"kubernetes_services_sidecar" json:"kubernetesServicesSidecar,omitempty"`
	// CoreDNS image
	CoreDNS string `yaml:"coredns" json:"coredns,omitempty"`
	// CoreDNS autoscaler image
	CoreDNSAutoscaler string `yaml:"coredns_autoscaler" json:"corednsAutoscaler,omitempty"`
	// Kubernetes image
	Kubernetes string `yaml:"kubernetes" json:"kubernetes,omitempty"`
	// Flannel image
	Flannel string `yaml:"flannel" json:"flannel,omitempty"`
	// Flannel CNI image
	FlannelCNI string `yaml:"flannel_cni" json:"flannelCni,omitempty"`
	// Calico Node image
	CalicoNode string `yaml:"calico_node" json:"calicoNode,omitempty"`
	// Calico CNI image
	CalicoCNI string `yaml:"calico_cni" json:"calicoCni,omitempty"`
	// Calico Controllers image
	CalicoControllers string `yaml:"calico_controllers" json:"calicoControllers,omitempty"`
	// Calicoctl image
	CalicoCtl string `yaml:"calico_ctl" json:"calicoCtl,omitempty"`
	// Pod infra container image
	PodInfraContainer string `yaml:"pod_infra_container" json:"podInfraContainer,omitempty"`
	// Ingress Controller image
	Ingress string `yaml:"ingress" json:"ingress,omitempty"`
	// Ingress Controller Backend image
	IngressBackend string `yaml:"ingress_backend" json:"ingressBackend,omitempty"`
	// Metrics Server image
	MetricsServer string `yaml:"metrics_server" json:"metricsServer,omitempty"`
	// Storage csi-attacher image
	StorageLvmAttacher        string `yaml:"storage_lvm_attacher" json:"storage_lvm_attacher,omitempty"`
	StorageLvmProvisioner     string `yaml:"storage_lvm_provisioner" json:"storage_lvm_provisioner,omitempty"`
	StorageLvmDriverRegistrar string `yaml:"storage_lvm_driver_registrar" json:"storage_lvm_driver_registrar,omitempty"`
	StorageLvmCSI             string `yaml:"storage_lvmcsi" json:"storage_lvmcsi,omitempty"`
	StorageLvmd               string `yaml:"storage_lvmd" json:"storage_lvmd,omitempty"`
	// Storage nfs image
	StorageNFSProvisioner      string `yaml:"storage_nfs" json:"storage_nfs,omitempty"`
	StorageNFSInit             string `yaml:"storage_nfs_init" json:"storage_nfs_init,omitempty"`
	ClusterAgent               string `yaml:"cluster_agent" json:"cluster_agent"`
	NodeAgent                  string `yaml:"node_agent" json:"node_agent"`
	StorageCephOperator        string `yaml:"storage_ceph_operator" json:"storage_ceph_operator,omitempty"`
	StorageCephCluster         string `yaml:"storage_ceph_cluster" json:"storage_ceph_cluster,omitempty"`
	StorageCephTools           string `yaml:"storage_ceph_tools" json:"storage_ceph_tools,omitempty"`
	StorageCephAttacher        string `yaml:"storage_ceph_attacher" json:"storage_ceph_attacher,omitempty"`
	StorageCephProvisioner     string `yaml:"storage_ceph_provisioner" json:"storage_ceph_provisioner,omitempty"`
	StorageCephDriverRegistrar string `yaml:"storage_ceph_driver_registrar" json:"storage_ceph_driver_registrar,omitempty"`
	StorageCephFsCSI           string `yaml:"storage_ceph_fscsi" json:"storage_ceph_fscsi,omitempty"`
	ZKERemover                 string `yaml:"zke_remover" json:zke_remover`
}

type ZKEConfigNode struct {
	// Name of the host provisioned via docker machine
	NodeName string `yaml:"nodeName,omitempty" json:"nodeName,omitempty" norman:"type=reference[node]"`
	// IP or FQDN that is fully resolvable and used for SSH communication
	Address string `yaml:"address" json:"address,omitempty"`
	// Port used for SSH communication
	Port string `yaml:"port" json:"port,omitempty"`
	// Optional - Internal address that will be used for components communication
	InternalAddress string `yaml:"internal_address" json:"internalAddress,omitempty"`
	// Node role in kubernetes cluster (controlplane, worker, or etcd)
	Role []string `yaml:"role" json:"role,omitempty" norman:"type=array[enum],options=etcd|worker|controlplane"`
	// Optional - Hostname of the node
	HostnameOverride string `yaml:"hostname_override" json:"hostnameOverride,omitempty"`
	// SSH usesr that will be used by ZKE
	User string `yaml:"user" json:"user,omitempty"`
	// Optional - Docker socket on the node that will be used in tunneling
	DockerSocket string `yaml:"docker_socket" json:"dockerSocket,omitempty"`
	// SSH Agent Auth enable
	SSHAgentAuth bool `yaml:"ssh_agent_auth,omitempty" json:"sshAgentAuth,omitempty"`
	// SSH Private Key
	SSHKey string `yaml:"ssh_key" json:"sshKey,omitempty" norman:"type=password"`
	// SSH Private Key Path
	SSHKeyPath string `yaml:"ssh_key_path" json:"sshKeyPath,omitempty"`
	// SSH Certificate
	SSHCert string `yaml:"ssh_cert" json:"sshCert,omitempty"`
	// SSH Certificate Path
	SSHCertPath string `yaml:"ssh_cert_path" json:"sshCertPath,omitempty"`
	// Node Labels
	Labels map[string]string `yaml:"labels" json:"labels,omitempty"`
}

type ZKEConfigServices struct {
	// Etcd Service
	Etcd ETCDService `yaml:"etcd" json:"etcd,omitempty"`
	// KubeAPI Service
	KubeAPI KubeAPIService `yaml:"kube-api" json:"kubeApi,omitempty"`
	// KubeController Service
	KubeController KubeControllerService `yaml:"kube-controller" json:"kubeController,omitempty"`
	// Scheduler Service
	Scheduler SchedulerService `yaml:"scheduler" json:"scheduler,omitempty"`
	// Kubelet Service
	Kubelet KubeletService `yaml:"kubelet" json:"kubelet,omitempty"`
	// KubeProxy Service
	Kubeproxy KubeproxyService `yaml:"kubeproxy" json:"kubeproxy,omitempty"`
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
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
}

type SchedulerService struct {
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
}

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

type NetworkConfig struct {
	// Network Plugin That will be used in kubernetes cluster
	Plugin string `yaml:"plugin" json:"plugin,omitempty" norman:"default=canal"`
	// Plugin options to configure network properties
	Options map[string]string `yaml:"options" json:"options,omitempty"`
	// CalicoNetworkProvider
	CalicoNetworkProvider *CalicoNetworkProvider `yaml:",omitempty" json:"calicoNetworkProvider,omitempty"`
	// FlannelNetworkProvider
	FlannelNetworkProvider *FlannelNetworkProvider `yaml:",omitempty" json:"flannelNetworkProvider,omitempty"`
}

type StorageConfig struct {
	Lvm  []Deviceconf `yaml:"lvm" json:"lvm,omitempty"`
	Nfs  []Deviceconf `yaml:"nfs" json:"nfs,omitempty"`
	Ceph []Deviceconf `yaml:"ceph" json:"ceph,omitempty"`
}

type Deviceconf struct {
	Host string   `yaml:"host" json:"host,omitempty"`
	Devs []string `yaml:"devs" json:"devs,omitempty"`
}

type AuthWebhookConfig struct {
	// ConfigFile is a multiline string that represent a custom webhook config file
	ConfigFile string `yaml:"config_file" json:"configFile,omitempty"`
	// CacheTimeout controls how long to cache authentication decisions
	CacheTimeout string `yaml:"cache_timeout" json:"cacheTimeout,omitempty"`
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

type IngressConfig struct {
	// Ingress controller type used by kubernetes
	Provider string `yaml:"provider" json:"provider,omitempty" norman:"default=nginx"`
	// Ingress controller options
	Options map[string]string `yaml:"options" json:"options,omitempty"`
	// NodeSelector key pair
	NodeSelector map[string]string `yaml:"node_selector" json:"nodeSelector,omitempty"`
	// Ingress controller extra arguments
	ExtraArgs map[string]string `yaml:"extra_args" json:"extraArgs,omitempty"`
}

type ZKEPlan struct {
	// List of node Plans
	Nodes []ZKEConfigNodePlan `json:"nodes,omitempty"`
}

type ZKEConfigNodePlan struct {
	// Node address
	Address string `json:"address,omitempty"`
	// map of named processes that should run on the node
	Processes map[string]Process `json:"processes,omitempty"`
	// List of portchecks that should be open on the node
	PortChecks []PortCheck `json:"portChecks,omitempty"`
	// List of files to deploy on the node
	Files []File `json:"files,omitempty"`
	// Node Annotations
	Annotations map[string]string `json:"annotations,omitempty"`
	// Node Labels
	Labels map[string]string `json:"labels,omitempty"`
}

type Process struct {
	// Process name, this should be the container name
	Name string `json:"name,omitempty"`
	// Process Entrypoint command
	Command []string `json:"command,omitempty"`
	// Process args
	Args []string `json:"args,omitempty"`
	// Environment variables list
	Env []string `json:"env,omitempty"`
	// Process docker image
	Image string `json:"image,omitempty"`
	//AuthConfig for image private registry
	ImageRegistryAuthConfig string `json:"imageRegistryAuthConfig,omitempty"`
	// Process docker image VolumesFrom
	VolumesFrom []string `json:"volumesFrom,omitempty"`
	// Process docker container bind mounts
	Binds []string `json:"binds,omitempty"`
	// Process docker container netwotk mode
	NetworkMode string `json:"networkMode,omitempty"`
	// Process container restart policy
	RestartPolicy string `json:"restartPolicy,omitempty"`
	// Process container pid mode
	PidMode string `json:"pidMode,omitempty"`
	// Run process in privileged container
	Privileged bool `json:"privileged,omitempty"`
	// Process healthcheck
	HealthCheck HealthCheck `json:"healthCheck,omitempty"`
	// Process docker container Labels
	Labels map[string]string `json:"labels,omitempty"`
	// Process docker publish container's port to host
	Publish []string `json:"publish,omitempty"`
}

type HealthCheck struct {
	// Healthcheck URL
	URL string `json:"url,omitempty"`
}

type PortCheck struct {
	// Portcheck address to check.
	Address string `json:"address,omitempty"`
	// Port number
	Port int `json:"port,omitempty"`
	// Port Protocol
	Protocol string `json:"protocol,omitempty"`
}

type CloudProvider struct {
	Name                string `yaml:"name" json:"name,omitempty"`
	CustomCloudProvider string `yaml:"customCloudProvider,omitempty" json:"customCloudProvider,omitempty"`
}

type CalicoNetworkProvider struct {
	CloudProvider string `json:"cloudProvider"`
}

type FlannelNetworkProvider struct {
	// Alternate cloud interface for flannel
	Iface string `json:"iface"`
}

type KubernetesServicesOptions struct {
	// Additional options passed to KubeAPI
	KubeAPI map[string]string `json:"kubeapi"`
	// Additional options passed to Kubelet
	Kubelet map[string]string `json:"kubelet"`
	// Additional options passed to Kubeproxy
	Kubeproxy map[string]string `json:"kubeproxy"`
	// Additional options passed to KubeController
	KubeController map[string]string `json:"kubeController"`
	// Additional options passed to Scheduler
	Scheduler map[string]string `json:"scheduler"`
}

type MonitorConfig struct {
	// Monitoring server provider
	MetricsProvider string `yaml:"metrics_provider" json:"metrics_provider,omitempty" norman:"default=metrics-server"`
	// Metrics server options
	MetricsOptions                        map[string]string `yaml:"metrics_options" json:"metrics_options,omitempty"`
	PrometheusAlertManagerIngressEndpoint string            `yaml:"prometheus_alertmanager_ingress_endpoint" json:"prometheus_alertmanager_ingress_endpoint"`
	GrafanaIngressEndpoint                string            `yaml:"grafana_ingress_endpoint" json:"grafana_ingress_endpoint"`
}

type RestoreConfig struct {
	Restore      bool   `yaml:"restore" json:"restore,omitempty"`
	SnapshotName string `yaml:"snapshot_name" json:"snapshotName,omitempty"`
}

type DNSConfig struct {
	// DNS provider
	Provider string `yaml:"provider" json:"provider,omitempty" norman:"coredns"`
	// Upstream nameservers
	UpstreamNameservers []string `yaml:"upstreamnameservers" json:"upstreamnameservers,omitempty"`
	// ReverseCIDRs
	ReverseCIDRs []string `yaml:"reversecidrs" json:"reversecidrs,omitempty"`
	// NodeSelector key pair
	NodeSelector map[string]string `yaml:"node_selector" json:"nodeSelector,omitempty"`
}

type File struct {
	Name     string `json:"name,omitempty"`
	Contents string `json:"contents,omitempty"`
}

type BackupConfig struct {
	// Backup interval in hours
	IntervalHours int `yaml:"interval_hours" json:"intervalHours,omitempty" norman:"default=12"`
	// Number of backups to keep
	Retention int `yaml:"retention" json:"retention,omitempty" norman:"default=6"`
}

type RegistryConfig struct {
	Isenabled               bool   `yaml:"isenabled" json:"isenabled,omitempty"`
	RegistryIngressURL      string `yaml:"registry_ingress_url" json:"registry_ingress_url"`
	NotaryIngressURL        string `yaml:"notary_ingress_url" json:"notary_ingress_url"`
	RegistryDiskCapacity    string `yaml:"registry_disk_capacity" json:"registry_disk_capacity"`
	DatabaseDiskCapacity    string `yaml:"database_disk_capacity" json:"database_disk_capacity"`
	RedisDiskCapacity       string `yaml:"redis_disk_capacity" json:"redis_disk_capacity"`
	ChartmuseumDiskCapacity string `yaml:"Chartmuseum_disk_capacity" json:"Chartmuseum_disk_capacity"`
	JobserviceDiskCapacity  string `yaml:"jobservice_disk_capacity" json:"jobservice_disk_capacity"`
}
