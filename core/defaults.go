package core

import (
	"context"
	"strings"

	"github.com/zdnscloud/zke/core/services"
	"github.com/zdnscloud/zke/pkg/docker"
	"github.com/zdnscloud/zke/pkg/k8s"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/types"
)

const (
	DefaultServiceClusterIPRange = "10.43.0.0/16"
	DefaultNodePortRange         = "30000-32767"
	DefaultClusterCIDR           = "10.42.0.0/16"
	DefaultClusterDNSService     = "10.43.0.10"
	DefaultClusterDomain         = "cluster.local"
	DefaultClusterName           = "local"
	DefaultClusterSSHKeyPath     = "~/.ssh/id_rsa"
	DefaultClusterGlobalDns      = "114.114.114.114,223.5.5.5"

	DefaultK8sVersion = types.DefaultK8s

	DefaultSSHPort        = "22"
	DefaultDockerSockPath = "/var/run/docker.sock"

	DefaultAuthStrategy      = "x509"
	DefaultAuthorizationMode = "rbac"

	DefaultAuthnWebhookFile = `
	apiVersion: v1
	kind: Config
	clusters:
	- name: Default
	  cluster:
		insecure-skip-tls-verify: true
		server: http://127.0.0.1:6440/v1/authenticate
	users:
	- name: Default
	  user:
		insecure-skip-tls-verify: true
	current-context: webhook
	contexts:
	- name: webhook
	  context:
		user: Default
		cluster: Default
	`
	DefaultAuthnCacheTimeout = "5s"

	DefaultNetworkPlugin      = "flannel"
	DefaultFlannelBackendType = "vxlan"

	DefaultNetworkCloudProvider = "none"

	DefaultDNSProvider = "coredns"

	DefaultStorageclass = "lvm"

	DefaultIngressController             = "nginx"
	DefaultEtcdBackupCreationPeriod      = "12h"
	DefaultEtcdBackupRetentionPeriod     = "72h"
	DefaultEtcdSnapshot                  = true
	DefaultEtcdBackupConfigIntervalHours = 12
	DefaultEtcdBackupConfigRetention     = 6

	DefaultEtcdHeartbeatIntervalName  = "heartbeat-interval"
	DefaultEtcdHeartbeatIntervalValue = "500"
	DefaultEtcdElectionTimeoutName    = "election-timeout"
	DefaultEtcdElectionTimeoutValue   = "5000"

	DefaultMonitoringMetricsProvider = "metrics-server"
	DefaultMonitoringNamespace       = "kube-monitoring"

	DefaultRegistryRedisDiskCapacity       = "1Gi"
	DefaultRegistryDatabaseDiskCapacity    = "5Gi"
	DefaultRegistryJobserviceDiskCapacity  = "1Gi"
	DefaultRegistryChartmuseumDiskCapacity = "5Gi"
)

type ExternalFlags struct {
	CertificateDir   string
	ClusterFilePath  string
	ConfigDir        string
	CustomCerts      bool
	DisablePortCheck bool
	GenerateCSR      bool
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
	if len(c.DNS.Provider) == 0 {
		c.DNS.Provider = DefaultDNSProvider
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
	if len(c.Monitor.MetricsProvider) == 0 {
		c.Monitor.MetricsProvider = DefaultMonitoringMetricsProvider
	}
	if len(c.Monitor.PrometheusAlertManagerIngressEndpoint) == 0 {
		c.Monitor.PrometheusAlertManagerIngressEndpoint = "alertmanager" + "." + DefaultMonitoringNamespace + "." + DefaultClusterDomain
	}
	if len(c.Monitor.GrafanaIngressEndpoint) == 0 {
		c.Monitor.GrafanaIngressEndpoint = "grafana" + "." + DefaultMonitoringNamespace + "." + DefaultClusterDomain
	}
	//set docker private registry URL
	for _, pr := range c.PrivateRegistries {
		if pr.URL == "" {
			pr.URL = docker.DockerRegistryURL
		}
		c.PrivateRegistriesMap[pr.URL] = pr
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

func (c *Cluster) setClusterNetworkDefaults() {
	setDefaultIfEmpty(&c.Network.Plugin, DefaultNetworkPlugin)

	if c.Network.Options == nil {
		// don't break if the user didn't define options
		c.Network.Options = make(map[string]string)
	}
	networkPluginConfigDefaultsMap := make(map[string]string)
	// This is still needed because ZKE doesn't use c.Network.*NetworkProvider
	switch c.Network.Plugin {
	case CalicoNetworkPlugin:
		networkPluginConfigDefaultsMap = map[string]string{
			CalicoCloudProvider: DefaultNetworkCloudProvider,
		}
	case FlannelNetworkPlugin:
		networkPluginConfigDefaultsMap = map[string]string{
			FlannelIface:                c.Network.Options[FlannelIface],
			FlannelBackendType:          c.Network.Options[FlannelBackendType],
			FlannelBackendDirectrouting: c.Network.Options[FlannelBackendDirectrouting],
		}
	}
	if c.Network.CalicoNetworkProvider != nil {
		setDefaultIfEmpty(&c.Network.CalicoNetworkProvider.CloudProvider, DefaultNetworkCloudProvider)
		networkPluginConfigDefaultsMap[CalicoCloudProvider] = c.Network.CalicoNetworkProvider.CloudProvider
	}
	if c.Network.FlannelNetworkProvider != nil {
		networkPluginConfigDefaultsMap[FlannelIface] = c.Network.FlannelNetworkProvider.Iface

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
		c.Authentication.Webhook = &types.AuthWebhookConfig{}
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

func GetExternalFlags(disablePortCheck bool, configDir, clusterFilePath string) ExternalFlags {
	return ExternalFlags{
		DisablePortCheck: disablePortCheck,
		ConfigDir:        configDir,
		ClusterFilePath:  clusterFilePath,
	}
}