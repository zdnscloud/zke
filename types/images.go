package types

type ZKESystemImages struct {
	Etcd string `yaml:"etcd" json:"etcd,omitempty"`
	// ZKE image
	Alpine                    string `yaml:"alpine" json:"alpine,omitempty"`
	NginxProxy                string `yaml:"nginx_proxy" json:"nginxProxy,omitempty"`
	CertDownloader            string `yaml:"cert_downloader" json:"certDownloader,omitempty"`
	ZKERemover                string `yaml:"zke_remover" json:zkeRemover`
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
	MetricsServer  string `yaml:"metrics_server" json:"metricsServer,omitempty"`
	// Zcloud image
	ClusterAgent string `yaml:"cluster_agent" json:"clusterAgent"`
	NodeAgent    string `yaml:"node_agent" json:"nodeAgent"`
}
