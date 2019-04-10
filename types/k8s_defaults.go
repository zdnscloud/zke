package types

import (
	"fmt"
	"strings"
)

const (
	DefaultK8s = "v1.13.1"
)

var (
	K8sBadVersions = map[string]bool{
		"v1.9.7":  true,
		"v1.10.1": true,
		"v1.8.11": true,
		"v1.8.10": true,
	}

	// K8sVersionsCurrent are the latest versions available for installation
	K8sVersionsCurrent = []string{
		"v1.13.1",
	}

	// K8sVersionToRKESystemImages is dynamically populated on init() with the latest versions
	K8sVersionToZKESystemImages map[string]ZKESystemImages

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

	m = func(image string) string {
		//orig := image
		if strings.HasPrefix(image, "weaveworks") {
			return image
		}
		image = strings.Replace(image, "gcr.io/google_containers", "zdnscloud", 1)
		image = strings.Replace(image, "quay.io/coreos/", "zdnscloud/coreos-", 1)
		image = strings.Replace(image, "quay.io/calico/", "zdnscloud/calico-", 1)
		image = strings.Replace(image, "k8s.gcr.io/", "zdnscloud/nginx-ingress-controller-", 1)
		image = strings.Replace(image, "coredns/", "zdnscloud/", 1)
		return image
	}

	AllK8sVersions = map[string]ZKESystemImages{
		"v1.13.1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.2.24"),
			Kubernetes:                m("zdnscloud/hyperkube:v1.13.1"),
			Alpine:                    m("zdnscloud/zke-tools:v0.1.23"),
			NginxProxy:                m("zdnscloud/zke-tools:v0.1.23"),
			CertDownloader:            m("zdnscloud/zke-tools:v0.1.23"),
			KubernetesServicesSidecar: m("zdnscloud/zke-tools:v0.1.23"),
			Flannel:                   m("quay.io/coreos/flannel:v0.10.0"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.3.0"),
			CalicoNode:                m("quay.io/calico/node:v3.4.0"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.4.0"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.1"),
			Ingress:                   m("zdnscloud/nginx-ingress-controller:0.21.0"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
			MetricsServer:             m("gcr.io/google_containers/metrics-server-amd64:v0.3.1"),
			CoreDNS:                   m("coredns/coredns:1.2.6"),
			CoreDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			StorageCSIAttacher:        m("quay.io/k8scsi/csi-attacher:v0.4.2"),
			StorageCSIProvisioner:     m("quay.io/k8scsi/csi-provisioner:v0.4.2"),
			StorageDriverRegistrar:    m("quay.io/k8scsi/driver-registrar:v0.4.2"),
			StorageCSILvmplugin:       m("quay.io/lvmcsi/lvmplugin:v0.3.1"),
			StorageLvmd:               m("zdnscloud/lvmd:v0.1"),
		},
	}
)

func init() {
	if K8sVersionToZKESystemImages != nil {
		panic("Do not initialize or add values to K8sVersionToRKESystemImages")
	}

	K8sVersionToZKESystemImages = map[string]ZKESystemImages{}

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
		K8sVersionToZKESystemImages[latest] = images
	}

	if _, ok := K8sVersionToZKESystemImages[DefaultK8s]; !ok {
		panic("Default K8s version " + DefaultK8s + " is not found in k8sVersionsCurrent list")
	}
}
