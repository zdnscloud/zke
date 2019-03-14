package v3

import (
	"fmt"
	"strings"

	projectv3 "github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/rancher/types/image"
)

const (
	DefaultK8s = "v1.13.1"
)

var (
	m = image.Mirror

	K8sBadVersions = map[string]bool{
		"v1.9.7-rancher1":    true,
		"v1.10.1-rancher1":   true,
		"v1.8.11-rancher1":   true,
		"v1.8.10-rancher1-1": true,
	}

	// K8sVersionsCurrent are the latest versions available for installation
	K8sVersionsCurrent = []string{
		"v1.13.1",
	}

	// K8sVersionToRKESystemImages is dynamically populated on init() with the latest versions
	K8sVersionToRKESystemImages map[string]RKESystemImages

	// K8sVersionServiceOptions - service options per k8s version
	K8sVersionServiceOptions = map[string]KubernetesServicesOptions{
		"v1.13": {
			KubeAPI: map[string]string{
				"tls-cipher-suites":        "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
				"enable-admission-plugins": "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota",
			},
			Kubelet: map[string]string{
				"tls-cipher-suites": "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
			},
		},
	}

	// ToolsSystemImages default images for alert, pipeline, logging
	ToolsSystemImages = struct {
		AlertSystemImages    AlertSystemImages
		PipelineSystemImages projectv3.PipelineSystemImages
		LoggingSystemImages  LoggingSystemImages
	}{
		AlertSystemImages: AlertSystemImages{
			AlertManager:       m("prom/alertmanager:v0.15.2"),
			AlertManagerHelper: m("rancher/alertmanager-helper:v0.0.2"),
		},
		PipelineSystemImages: projectv3.PipelineSystemImages{
			Jenkins:       m("rancher/pipeline-jenkins-server:v0.1.0"),
			JenkinsJnlp:   m("jenkins/jnlp-slave:3.10-1-alpine"),
			AlpineGit:     m("rancher/pipeline-tools:v0.1.4"),
			PluginsDocker: m("plugins/docker:17.12"),
			Minio:         m("minio/minio:RELEASE.2018-05-25T19-49-13Z"),
			Registry:      m("registry:2"),
			RegistryProxy: m("rancher/pipeline-tools:v0.1.4"),
			KubeApply:     m("rancher/pipeline-tools:v0.1.4"),
		},
		LoggingSystemImages: LoggingSystemImages{
			Fluentd:                       m("rancher/fluentd:v0.1.11"),
			FluentdHelper:                 m("rancher/fluentd-helper:v0.1.2"),
			LogAggregatorFlexVolumeDriver: m("rancher/log-aggregator:v0.1.4"),
		},
	}

	AllK8sVersions = map[string]RKESystemImages{
		"v1.13.1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.2.24"),
			Kubernetes:                m("zdnscloud/hyperkube:v1.13.1"),
			Alpine:                    m("zdnscloud/zke-tools:v0.1.23"),
			NginxProxy:                m("zdnscloud/zke-tools:v0.1.23"),
			CertDownloader:            m("zdnscloud/zke-tools:v0.1.23"),
			KubernetesServicesSidecar: m("zdnscloud/zke-tools:v0.1.23"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.15.0"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.15.0"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.15.0"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.10.0"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.3.0"),
			CalicoNode:                m("quay.io/calico/node:v3.4.0"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.4.0"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.4.0"),
			CanalCNI:                  m("quay.io/calico/cni:v3.4.0"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.10.0"),
			WeaveNode:                 m("weaveworks/weave-kube:2.5.0"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.5.0"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.1"),
			Ingress:                   m("zdnscloud/nginx-ingress-controller:0.21.0"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
			MetricsServer:             m("gcr.io/google_containers/metrics-server-amd64:v0.3.1"),
			CoreDNS:                   m("coredns/coredns:1.2.6"),
			CoreDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
		},
	}
)

func init() {
	if K8sVersionToRKESystemImages != nil {
		panic("Do not initialize or add values to K8sVersionToRKESystemImages")
	}

	K8sVersionToRKESystemImages = map[string]RKESystemImages{}

	for version, images := range AllK8sVersions {
		if K8sBadVersions[version] {
			continue
		}

		longName := "zdnscloud/hyperkube:" + version
		if !strings.HasPrefix(longName, images.Kubernetes) {
			panic(fmt.Sprintf("For K8s version %s, the Kubernetes image tag should be a substring of %s, currently it is %s", version, version, images.Kubernetes))
		}
	}

	for _, latest := range K8sVersionsCurrent {
		images, ok := AllK8sVersions[latest]
		if !ok {
			panic("K8s version " + " is not found in AllK8sVersions map")
		}
		K8sVersionToRKESystemImages[latest] = images
	}

	if _, ok := K8sVersionToRKESystemImages[DefaultK8s]; !ok {
		panic("Default K8s version " + DefaultK8s + " is not found in k8sVersionsCurrent list")
	}

	// init Windows versions
	initWindows()
}
