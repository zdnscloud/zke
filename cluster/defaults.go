package cluster

import (
	"context"
	"fmt"
	"strings"

	"github.com/rancher/rke/cloudprovider"
	"github.com/rancher/rke/docker"
	"github.com/rancher/rke/k8s"
	"github.com/rancher/rke/log"
	"github.com/rancher/rke/services"
	"github.com/rancher/rke/templates"
	"github.com/rancher/rke/util"
	"github.com/rancher/types/apis/management.cattle.io/v3"
)

const (
	DefaultServiceClusterIPRange = "10.43.0.0/16"
	DefaultNodePortRange         = "30000-32767"
	DefaultClusterCIDR           = "10.42.0.0/16"
	DefaultClusterDNSService     = "10.43.0.10"
	DefaultClusterDomain         = "cluster.local"
	DefaultClusterName           = "local"
	DefaultClusterSSHKeyPath     = "~/.ssh/id_rsa"

	DefaultK8sVersion = v3.DefaultK8s

	DefaultSSHPort        = "22"
	DefaultDockerSockPath = "/var/run/docker.sock"

	DefaultAuthStrategy      = "x509"
	DefaultAuthorizationMode = "rbac"

	DefaultAuthnWebhookFile  = templates.AuthnWebhook
	DefaultAuthnCacheTimeout = "5s"

	DefaultNetworkPlugin        = "canal"
	DefaultNetworkCloudProvider = "none"

	DefaultIngressController             = "nginx"
	DefaultEtcdBackupCreationPeriod      = "12h"
	DefaultEtcdBackupRetentionPeriod     = "72h"
	DefaultEtcdSnapshot                  = true
	DefaultMonitoringProvider            = "metrics-server"
	DefaultEtcdBackupConfigIntervalHours = 12
	DefaultEtcdBackupConfigRetention     = 6
	DefaultDNSProvider                   = "kubedns"

	DefaultEtcdHeartbeatIntervalName  = "heartbeat-interval"
	DefaultEtcdHeartbeatIntervalValue = "500"
	DefaultEtcdElectionTimeoutName    = "election-timeout"
	DefaultEtcdElectionTimeoutValue   = "5000"
)

type ExternalFlags struct {
	CertificateDir   string
	ClusterFilePath  string
	DinD             bool
	ConfigDir        string
	CustomCerts      bool
	DisablePortCheck bool
	GenerateCSR      bool
	Local            bool
	UpdateOnly       bool
}

func setDefaultIfEmptyMapValue(configMap map[string]string, key string, value string) {
	if _, ok := configMap[key]; !ok {
		configMap[key] = value
	}
}

func setDefaultIfEmpty(varName *string, defaultValue string) {
	if len(*varName) == 0 {
		*varName = defaultValue
	}
}

func (c *Cluster) setClusterDefaults(ctx context.Context) error {
	if len(c.SSHKeyPath) == 0 {
		c.SSHKeyPath = DefaultClusterSSHKeyPath
	}
	// Default Path prefix
	if len(c.PrefixPath) == 0 {
		c.PrefixPath = "/"
	}
	// Set bastion/jump host defaults
	if len(c.BastionHost.Address) > 0 {
		if len(c.BastionHost.Port) == 0 {
			c.BastionHost.Port = DefaultSSHPort
		}
		if len(c.BastionHost.SSHKeyPath) == 0 {
			c.BastionHost.SSHKeyPath = c.SSHKeyPath
		}
		c.BastionHost.SSHAgentAuth = c.SSHAgentAuth

	}
	for i, host := range c.Nodes {
		if len(host.InternalAddress) == 0 {
			c.Nodes[i].InternalAddress = c.Nodes[i].Address
		}
		if len(host.HostnameOverride) == 0 {
			// This is a temporary modification
			c.Nodes[i].HostnameOverride = c.Nodes[i].Address
		}
		if len(host.SSHKeyPath) == 0 {
			c.Nodes[i].SSHKeyPath = c.SSHKeyPath
		}
		if len(host.Port) == 0 {
			c.Nodes[i].Port = DefaultSSHPort
		}

		c.Nodes[i].HostnameOverride = strings.ToLower(c.Nodes[i].HostnameOverride)
		// For now, you can set at the global level only.
		c.Nodes[i].SSHAgentAuth = c.SSHAgentAuth
	}

	if len(c.Authorization.Mode) == 0 {
		c.Authorization.Mode = DefaultAuthorizationMode
	}
	if c.Services.KubeAPI.PodSecurityPolicy && c.Authorization.Mode != services.RBACAuthorizationMode {
		log.Warnf(ctx, "PodSecurityPolicy can't be enabled with RBAC support disabled")
		c.Services.KubeAPI.PodSecurityPolicy = false
	}
	if len(c.Ingress.Provider) == 0 {
		c.Ingress.Provider = DefaultIngressController
	}
	if len(c.ClusterName) == 0 {
		c.ClusterName = DefaultClusterName
	}
	if len(c.Version) == 0 {
		c.Version = DefaultK8sVersion
	}
	if c.AddonJobTimeout == 0 {
		c.AddonJobTimeout = k8s.DefaultTimeout
	}
	if len(c.Monitoring.Provider) == 0 {
		c.Monitoring.Provider = DefaultMonitoringProvider
	}
	//set docker private registry URL
	for _, pr := range c.PrivateRegistries {
		if pr.URL == "" {
			pr.URL = docker.DockerRegistryURL
		}
		c.PrivateRegistriesMap[pr.URL] = pr
	}

	err := c.setClusterImageDefaults()
	if err != nil {
		return err
	}

	if len(c.DNS.Provider) == 0 {
		c.DNS.Provider = DefaultDNSProvider
	}

	c.setClusterServicesDefaults()
	c.setClusterNetworkDefaults()
	c.setClusterAuthnDefaults()

	return nil
}

func (c *Cluster) setClusterServicesDefaults() {
	// We don't accept per service images anymore.
	c.Services.KubeAPI.Image = c.SystemImages.Kubernetes
	c.Services.Scheduler.Image = c.SystemImages.Kubernetes
	c.Services.KubeController.Image = c.SystemImages.Kubernetes
	c.Services.Kubelet.Image = c.SystemImages.Kubernetes
	c.Services.Kubeproxy.Image = c.SystemImages.Kubernetes
	c.Services.Etcd.Image = c.SystemImages.Etcd

	// enable etcd snapshots by default
	if c.Services.Etcd.Snapshot == nil {
		defaultSnapshot := DefaultEtcdSnapshot
		c.Services.Etcd.Snapshot = &defaultSnapshot
	}

	serviceConfigDefaultsMap := map[*string]string{
		&c.Services.KubeAPI.ServiceClusterIPRange:        DefaultServiceClusterIPRange,
		&c.Services.KubeAPI.ServiceNodePortRange:         DefaultNodePortRange,
		&c.Services.KubeController.ServiceClusterIPRange: DefaultServiceClusterIPRange,
		&c.Services.KubeController.ClusterCIDR:           DefaultClusterCIDR,
		&c.Services.Kubelet.ClusterDNSServer:             DefaultClusterDNSService,
		&c.Services.Kubelet.ClusterDomain:                DefaultClusterDomain,
		&c.Services.Kubelet.InfraContainerImage:          c.SystemImages.PodInfraContainer,
		&c.Services.Etcd.Creation:                        DefaultEtcdBackupCreationPeriod,
		&c.Services.Etcd.Retention:                       DefaultEtcdBackupRetentionPeriod,
	}
	for k, v := range serviceConfigDefaultsMap {
		setDefaultIfEmpty(k, v)
	}
	// Add etcd timeouts
	if c.Services.Etcd.ExtraArgs == nil {
		c.Services.Etcd.ExtraArgs = make(map[string]string)
	}
	if _, ok := c.Services.Etcd.ExtraArgs[DefaultEtcdElectionTimeoutName]; !ok {
		c.Services.Etcd.ExtraArgs[DefaultEtcdElectionTimeoutName] = DefaultEtcdElectionTimeoutValue
	}
	if _, ok := c.Services.Etcd.ExtraArgs[DefaultEtcdHeartbeatIntervalName]; !ok {
		c.Services.Etcd.ExtraArgs[DefaultEtcdHeartbeatIntervalName] = DefaultEtcdHeartbeatIntervalValue
	}

	if c.Services.Etcd.BackupConfig != nil {
		if c.Services.Etcd.BackupConfig.IntervalHours == 0 {
			c.Services.Etcd.BackupConfig.IntervalHours = DefaultEtcdBackupConfigIntervalHours
		}
		if c.Services.Etcd.BackupConfig.Retention == 0 {
			c.Services.Etcd.BackupConfig.Retention = DefaultEtcdBackupConfigRetention
		}
	}
}

func (c *Cluster) setClusterImageDefaults() error {
	var privRegURL string

	// Version Check
	err := util.ValidateVersion(c.Version)
	if err != nil {
		return err
	}

	imageDefaults := v3.AllK8sVersions[c.Version]

	for _, privReg := range c.PrivateRegistries {
		if privReg.IsDefault {
			privRegURL = privReg.URL
			break
		}
	}
	systemImagesDefaultsMap := map[*string]string{
		&c.SystemImages.Alpine:                    d(imageDefaults.Alpine, privRegURL),
		&c.SystemImages.NginxProxy:                d(imageDefaults.NginxProxy, privRegURL),
		&c.SystemImages.CertDownloader:            d(imageDefaults.CertDownloader, privRegURL),
		&c.SystemImages.KubeDNS:                   d(imageDefaults.KubeDNS, privRegURL),
		&c.SystemImages.KubeDNSSidecar:            d(imageDefaults.KubeDNSSidecar, privRegURL),
		&c.SystemImages.DNSmasq:                   d(imageDefaults.DNSmasq, privRegURL),
		&c.SystemImages.KubeDNSAutoscaler:         d(imageDefaults.KubeDNSAutoscaler, privRegURL),
		&c.SystemImages.CoreDNS:                   d(imageDefaults.CoreDNS, privRegURL),
		&c.SystemImages.CoreDNSAutoscaler:         d(imageDefaults.CoreDNSAutoscaler, privRegURL),
		&c.SystemImages.KubernetesServicesSidecar: d(imageDefaults.KubernetesServicesSidecar, privRegURL),
		&c.SystemImages.Etcd:                      d(imageDefaults.Etcd, privRegURL),
		&c.SystemImages.Kubernetes:                d(imageDefaults.Kubernetes, privRegURL),
		&c.SystemImages.PodInfraContainer:         d(imageDefaults.PodInfraContainer, privRegURL),
		&c.SystemImages.Flannel:                   d(imageDefaults.Flannel, privRegURL),
		&c.SystemImages.FlannelCNI:                d(imageDefaults.FlannelCNI, privRegURL),
		&c.SystemImages.CalicoNode:                d(imageDefaults.CalicoNode, privRegURL),
		&c.SystemImages.CalicoCNI:                 d(imageDefaults.CalicoCNI, privRegURL),
		&c.SystemImages.CalicoCtl:                 d(imageDefaults.CalicoCtl, privRegURL),
		&c.SystemImages.CanalNode:                 d(imageDefaults.CanalNode, privRegURL),
		&c.SystemImages.CanalCNI:                  d(imageDefaults.CanalCNI, privRegURL),
		&c.SystemImages.CanalFlannel:              d(imageDefaults.CanalFlannel, privRegURL),
		&c.SystemImages.WeaveNode:                 d(imageDefaults.WeaveNode, privRegURL),
		&c.SystemImages.WeaveCNI:                  d(imageDefaults.WeaveCNI, privRegURL),
		&c.SystemImages.Ingress:                   d(imageDefaults.Ingress, privRegURL),
		&c.SystemImages.IngressBackend:            d(imageDefaults.IngressBackend, privRegURL),
		&c.SystemImages.MetricsServer:             d(imageDefaults.MetricsServer, privRegURL),
	}

	for k, v := range systemImagesDefaultsMap {
		setDefaultIfEmpty(k, v)
	}

	return nil
}

func (c *Cluster) setClusterNetworkDefaults() {
	setDefaultIfEmpty(&c.Network.Plugin, DefaultNetworkPlugin)

	if c.Network.Options == nil {
		// don't break if the user didn't define options
		c.Network.Options = make(map[string]string)
	}
	networkPluginConfigDefaultsMap := make(map[string]string)
	// This is still needed because RKE doesn't use c.Network.*NetworkProvider, that's a rancher type
	switch c.Network.Plugin {
	case CalicoNetworkPlugin:
		networkPluginConfigDefaultsMap = map[string]string{
			CalicoCloudProvider: DefaultNetworkCloudProvider,
		}
	case FlannelNetworkPlugin:
		networkPluginConfigDefaultsMap = map[string]string{
			FlannelBackendType: "vxlan",
		}
	case CanalNetworkPlugin:
		networkPluginConfigDefaultsMap = map[string]string{
			CanalFlannelBackendType: "vxlan",
		}
	}
	if c.Network.CalicoNetworkProvider != nil {
		setDefaultIfEmpty(&c.Network.CalicoNetworkProvider.CloudProvider, DefaultNetworkCloudProvider)
		networkPluginConfigDefaultsMap[CalicoCloudProvider] = c.Network.CalicoNetworkProvider.CloudProvider
	}
	if c.Network.FlannelNetworkProvider != nil {
		networkPluginConfigDefaultsMap[FlannelIface] = c.Network.FlannelNetworkProvider.Iface

	}
	if c.Network.CanalNetworkProvider != nil {
		networkPluginConfigDefaultsMap[CanalIface] = c.Network.CanalNetworkProvider.Iface
	}
	if c.Network.WeaveNetworkProvider != nil {
		networkPluginConfigDefaultsMap[WeavePassword] = c.Network.WeaveNetworkProvider.Password
	}
	for k, v := range networkPluginConfigDefaultsMap {
		setDefaultIfEmptyMapValue(c.Network.Options, k, v)
	}
}

func (c *Cluster) setClusterAuthnDefaults() {
	setDefaultIfEmpty(&c.Authentication.Strategy, DefaultAuthStrategy)

	for _, strategy := range strings.Split(c.Authentication.Strategy, "|") {
		strategy = strings.ToLower(strings.TrimSpace(strategy))
		c.AuthnStrategies[strategy] = true
	}

	if c.AuthnStrategies[AuthnWebhookProvider] && c.Authentication.Webhook == nil {
		c.Authentication.Webhook = &v3.AuthWebhookConfig{}
	}
	if c.Authentication.Webhook != nil {
		webhookConfigDefaultsMap := map[*string]string{
			&c.Authentication.Webhook.ConfigFile:   DefaultAuthnWebhookFile,
			&c.Authentication.Webhook.CacheTimeout: DefaultAuthnCacheTimeout,
		}
		for k, v := range webhookConfigDefaultsMap {
			setDefaultIfEmpty(k, v)
		}
	}
}

func d(image, defaultRegistryURL string) string {
	if len(defaultRegistryURL) == 0 {
		return image
	}
	return fmt.Sprintf("%s/%s", defaultRegistryURL, image)
}

func (c *Cluster) setCloudProvider() error {
	p, err := cloudprovider.InitCloudProvider(c.CloudProvider)
	if err != nil {
		return fmt.Errorf("Failed to initialize cloud provider: %v", err)
	}
	if p != nil {
		c.CloudConfigFile, err = p.GenerateCloudConfigFile()
		if err != nil {
			return fmt.Errorf("Failed to parse cloud config file: %v", err)
		}
		c.CloudProvider.Name = p.GetName()
		if c.CloudProvider.Name == "" {
			return fmt.Errorf("Name of the cloud provider is not defined for custom provider")
		}
	}
	return nil
}

func GetExternalFlags(local, updateOnly, disablePortCheck bool, configDir, clusterFilePath string) ExternalFlags {
	return ExternalFlags{
		Local:            local,
		UpdateOnly:       updateOnly,
		DisablePortCheck: disablePortCheck,
		ConfigDir:        configDir,
		ClusterFilePath:  clusterFilePath,
	}
}
