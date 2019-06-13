package types

type ZKESystemImages struct {
	Etcd string `yaml:"etcd" json:"etcd,omitempty"`
	// ZKE image
	Alpine                    string `yaml:"alpine" json:"alpine,omitempty"`
	NginxProxy                string `yaml:"nginx_proxy" json:"nginxProxy,omitempty"`
	CertDownloader            string `yaml:"cert_downloader" json:"certDownloader,omitempty"`
	ZKERemover                string `yaml:"zke_remover" json:zke_remover`
	KubernetesServicesSidecar string `yaml:"kubernetes_services_sidecar" json:"kubernetesServicesSidecar,omitempty"`
	// CoreDNS image
	CoreDNS           string `yaml:"coredns" json:"coredns,omitempty"`
	CoreDNSAutoscaler string `yaml:"coredns_autoscaler" json:"corednsAutoscaler,omitempty"`
	// Kubernetes image
	Kubernetes string `yaml:"kubernetes" json:"kubernetes,omitempty"`
	// Flannel image
	Flannel    string `yaml:"flannel" json:"flannel,omitempty"`
	FlannelCNI string `yaml:"flannel_cni" json:"flannelCni,omitempty"`
	// Calico image
	CalicoNode        string `yaml:"calico_node" json:"calicoNode,omitempty"`
	CalicoCNI         string `yaml:"calico_cni" json:"calicoCni,omitempty"`
	CalicoControllers string `yaml:"calico_controllers" json:"calicoControllers,omitempty"`
	CalicoCtl         string `yaml:"calico_ctl" json:"calicoCtl,omitempty"`
	// Pod infra container image
	PodInfraContainer string `yaml:"pod_infra_container" json:"podInfraContainer,omitempty"`
	// Ingress Controller image
	Ingress        string `yaml:"ingress" json:"ingress,omitempty"`
	IngressBackend string `yaml:"ingress_backend" json:"ingressBackend,omitempty"`
	// Storage csi-attacher image
	StorageLvmAttacher        string `yaml:"storage_lvm_attacher" json:"storage_lvm_attacher,omitempty"`
	StorageLvmProvisioner     string `yaml:"storage_lvm_provisioner" json:"storage_lvm_provisioner,omitempty"`
	StorageLvmDriverRegistrar string `yaml:"storage_lvm_driver_registrar" json:"storage_lvm_driver_registrar,omitempty"`
	StorageLvmCSI             string `yaml:"storage_lvmcsi" json:"storage_lvmcsi,omitempty"`
	StorageLvmd               string `yaml:"storage_lvmd" json:"storage_lvmd,omitempty"`
	// Storage nfs image
	StorageNFSProvisioner      string `yaml:"storage_nfs" json:"storage_nfs,omitempty"`
	StorageNFSInit             string `yaml:"storage_nfs_init" json:"storage_nfs_init,omitempty"`
	StorageCephOperator        string `yaml:"storage_ceph_operator" json:"storage_ceph_operator,omitempty"`
	StorageCephCluster         string `yaml:"storage_ceph_cluster" json:"storage_ceph_cluster,omitempty"`
	StorageCephTools           string `yaml:"storage_ceph_tools" json:"storage_ceph_tools,omitempty"`
	StorageCephAttacher        string `yaml:"storage_ceph_attacher" json:"storage_ceph_attacher,omitempty"`
	StorageCephProvisioner     string `yaml:"storage_ceph_provisioner" json:"storage_ceph_provisioner,omitempty"`
	StorageCephDriverRegistrar string `yaml:"storage_ceph_driver_registrar" json:"storage_ceph_driver_registrar,omitempty"`
	StorageCephFsCSI           string `yaml:"storage_ceph_fscsi" json:"storage_ceph_fscsi,omitempty"`
	// Harbor Registry image
	HarborAdminserver  string `yaml:"harbor_adminserver" json:"harbor_adminserver"`
	HarborChartmuseum  string `yaml:"harbor_chartmuseum" json:"harbor_chartmuseum"`
	HarborClair        string `yaml:"harbor_clair" json:"harbor_clair"`
	HarborCore         string `yaml:"harbor_core" json:"harbor_core"`
	HarborDatabase     string `yaml:"harbor_database" json:"harbor_database"`
	HarborJobservice   string `yaml:"harbor_jobservice" json:"harbor_jobservice"`
	HarborNotaryServer string `yaml:"harbor_notaryserver" json:"harbor_notaryserver"`
	HarborNotarySigner string `yaml:"harbor_notarysigner" json:"harbor_notarysigner"`
	HarborPortal       string `yaml:"harbor_portal" json:"harbor_portal"`
	HarborRedis        string `yaml:"harbor_redis" json:"harbor_redis"`
	HarborRegistry     string `yaml:"harbor_registry" json:"harbor_registry"`
	HarborRegistryctl  string `yaml:"harbor_registryctl" json:"harbor_registryctl"`
	// Monitor image
	PrometheusAlertManager      string `yaml:"prometheus_alert_manager" json:"prometheus_alert_manager"`
	PrometheusConfigMapReloader string `yaml:"prometheus_configmap_reloader" json:"prometheus_configmap_reloader"`
	PrometheusNodeExporter      string `yaml:"prometheus_nodeexporter" json:"prometheus_nodeexporter"`
	PrometheusServer            string `yaml:"prometheus_server" json:"prometheus_server"`
	Grafana                     string `yaml:"grafana" json:"grafana"`
	GrafanaWatcher              string `yaml:"grafana_watcher" json:"grafana_watcher"`
	KubeStateMetrics            string `yaml:"kube_state_metrics" json:"kube_state_metrics"`
	MetricsServer               string `yaml:"metrics_server" json:"metricsServer,omitempty"`
	// Zcloud image
	ClusterAgent string `yaml:"cluster_agent" json:"cluster_agent"`
	NodeAgent    string `yaml:"node_agent" json:"node_agent"`
}
