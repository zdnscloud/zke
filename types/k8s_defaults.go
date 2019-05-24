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

	// K8sVersionToZKESystemImages is dynamically populated on init() with the latest versions
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
			Etcd:                       m("quay.io/coreos/etcd:v3.2.24"),
			Kubernetes:                 m("zdnscloud/hyperkube:v1.13.1"),
			Alpine:                     m("zdnscloud/zke-tools:v0.1.36"),
			NginxProxy:                 m("zdnscloud/zke-tools:v0.1.36"),
			CertDownloader:             m("zdnscloud/zke-tools:v0.1.36"),
			KubernetesServicesSidecar:  m("zdnscloud/zke-tools:v0.1.36"),
			Flannel:                    m("quay.io/coreos/flannel:v0.10.0"),
			FlannelCNI:                 m("quay.io/coreos/flannel-cni:v0.3.0"),
			CalicoNode:                 m("quay.io/calico/node:v3.4.0"),
			CalicoCNI:                  m("quay.io/calico/cni:v3.4.0"),
			CalicoCtl:                  m("quay.io/calico/ctl:v2.0.0"),
			PodInfraContainer:          m("gcr.io/google_containers/pause-amd64:3.1"),
			Ingress:                    m("zdnscloud/nginx-ingress-controller:0.21.0"),
			IngressBackend:             m("k8s.gcr.io/defaultbackend:1.4"),
			CoreDNS:                    m("coredns/coredns:1.2.6"),
			CoreDNSAutoscaler:          m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			StorageLvmAttacher:         m("quay.io/k8scsi/csi-attacher:v1.0.0"),
			StorageLvmProvisioner:      m("quay.io/k8scsi/csi-provisioner:v1.0.0"),
			StorageLvmDriverRegistrar:  m("quay.io/k8scsi/csi-node-driver-registrar:v1.0.2"),
			StorageLvmCSI:              m("zdnscloud/lvmcsi:v0.3"),
			StorageLvmd:                m("zdnscloud/lvmd:v0.3"),
			StorageNFSProvisioner:      m("quay.io/kubernetes_incubator/nfs-provisioner:v2.2.1-k8s1.12"),
			StorageNFSInit:             m("zdnscloud/nfs-init:v0.1"),
			ClusterAgent:               m("zdnscloud/cluster-agent:v1.0"),
			NodeAgent:                  m("zdnscloud/node-agent:v1.0"),
			StorageCephOperator:        m("rook/ceph:master"),
			StorageCephCluster:         m("ceph/ceph:v14.2.1-20190430"),
			StorageCephTools:           m("rook/ceph:master"),
			StorageCephAttacher:        m("quay.io/k8scsi/csi-attacher:v1.0.1"),
			StorageCephProvisioner:     m("quay.io/k8scsi/csi-provisioner:v1.0.1"),
			StorageCephDriverRegistrar: m("quay.io/k8scsi/csi-node-driver-registrar:v1.0.2"),
			StorageCephFsCSI:           m("quay.io/cephcsi/cephfsplugin:v1.0.0"),
		},
	}
)

func init() {
	if K8sVersionToZKESystemImages != nil {
		panic("Do not initialize or add values to K8sVersionToZKESystemImages")
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
