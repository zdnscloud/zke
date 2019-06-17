package typesnew

type ZKESystemImages struct {
	Etcd string `json:"etcd"`
	// ZKE image
	Alpine                    string `json:"alpine"`
	NginxProxy                string `json:"nginxProxy"`
	CertDownloader            string `json:"certDownloader"`
	ZKERemover                string `json:zkeRemover`
	KubernetesServicesSidecar string `json:"kubernetesServicesSidecar"`
	// CoreDNS image
	CoreDNS           string `json:"coredns"`
	CoreDNSAutoscaler string `json:"corednsAutoscaler"`
	// Kubernetes image
	Kubernetes string `json:"kubernetes"`
	// Flannel image
	Flannel    string `json:"flannel"`
	FlannelCNI string `json:"flannelCni"`
	// Calico image
	CalicoNode        string `json:"calicoNode"`
	CalicoCNI         string `json:"calicoCni"`
	CalicoControllers string `json:"calicoControllers"`
	CalicoCtl         string `json:"calicoCtl"`
	// Pod infra container image
	PodInfraContainer string `json:"podInfraContainer"`
	// Ingress Controller image
	Ingress        string `json:"ingress"`
	IngressBackend string `json:"ingressBackend"`
	MetricsServer  string `json:"metricsServer"`
	// Zcloud image
	ClusterAgent string `json:"clusterAgent"`
	NodeAgent    string `json:"nodeAgent"`
}
