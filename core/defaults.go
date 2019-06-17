package core

import (
	"context"
	"strings"

	"github.com/zdnscloud/zke/core/services"
	"github.com/zdnscloud/zke/pkg/docker"
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
	if len(c.Option.SSHKeyPath) == 0 {
		c.Option.SSHKeyPath = DefaultClusterSSHKeyPath
	}
	// Default Path prefix
	if len(c.Option.PrefixPath) == 0 {
		c.Option.PrefixPath = "/"
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
			c.Nodes[i].SSHKeyPath = c.Option.SSHKeyPath
		}
		if len(host.Port) == 0 {
			c.Nodes[i].Port = DefaultSSHPort
		}

		c.Nodes[i].HostnameOverride = strings.ToLower(c.Nodes[i].HostnameOverride)
	}

	if len(c.Authorization.Mode) == 0 {
		c.Authorization.Mode = DefaultAuthorizationMode
	}
	if c.Core.KubeAPI.PodSecurityPolicy && c.Authorization.Mode != services.RBACAuthorizationMode {
		log.Warnf(ctx, "PodSecurityPolicy can't be enabled with RBAC support disabled")
		c.Core.KubeAPI.PodSecurityPolicy = false
	}
	if len(c.Network.Ingress.Provider) == 0 {
		c.Network.Ingress.Provider = DefaultIngressController
	}
	if len(c.Network.DNS.Provider) == 0 {
		c.Network.DNS.Provider = DefaultDNSProvider
	}
	if len(c.ClusterName) == 0 {
		c.ClusterName = DefaultClusterName
	}
	if len(c.Option.KubernetesVersion) == 0 {
		c.Option.KubernetesVersion = DefaultK8sVersion
	}
	if len(c.Monitor.MetricsProvider) == 0 {
		c.Monitor.MetricsProvider = DefaultMonitoringMetricsProvider
	}
	//set docker private registry URL
	for _, pr := range c.PrivateRegistries {
		if pr.URL == "" {
			pr.URL = docker.DockerRegistryURL
		}
		c.PrivateRegistriesMap[pr.URL] = pr
	}

	c.setClusterServicesDefaults()
	c.setClusterAuthnDefaults()

	return nil
}

func (c *Cluster) setClusterServicesDefaults() {
	// We don't accept per service images anymore.
	c.Core.KubeAPI.Image = c.SystemImages.Kubernetes
	c.Core.Scheduler.Image = c.SystemImages.Kubernetes
	c.Core.KubeController.Image = c.SystemImages.Kubernetes
	c.Core.Kubelet.Image = c.SystemImages.Kubernetes
	c.Core.Kubeproxy.Image = c.SystemImages.Kubernetes
	c.Core.Etcd.Image = c.SystemImages.Etcd

	// enable etcd snapshots by default
	if c.Core.Etcd.Snapshot == nil {
		defaultSnapshot := DefaultEtcdSnapshot
		c.Core.Etcd.Snapshot = &defaultSnapshot
	}

	serviceConfigDefaultsMap := map[*string]string{
		&c.Core.KubeAPI.ServiceClusterIPRange:        DefaultServiceClusterIPRange,
		&c.Core.KubeAPI.ServiceNodePortRange:         DefaultNodePortRange,
		&c.Core.KubeController.ServiceClusterIPRange: DefaultServiceClusterIPRange,
		&c.Core.KubeController.ClusterCIDR:           DefaultClusterCIDR,
		&c.Core.Kubelet.ClusterDNSServer:             DefaultClusterDNSService,
		&c.Core.Kubelet.ClusterDomain:                DefaultClusterDomain,
		&c.Core.Kubelet.InfraContainerImage:          c.SystemImages.PodInfraContainer,
		&c.Core.Etcd.Creation:                        DefaultEtcdBackupCreationPeriod,
		&c.Core.Etcd.Retention:                       DefaultEtcdBackupRetentionPeriod,
	}
	for k, v := range serviceConfigDefaultsMap {
		setDefaultIfEmpty(k, v)
	}
	// Add etcd timeouts
	if c.Core.Etcd.ExtraArgs == nil {
		c.Core.Etcd.ExtraArgs = make(map[string]string)
	}
	if _, ok := c.Core.Etcd.ExtraArgs[DefaultEtcdElectionTimeoutName]; !ok {
		c.Core.Etcd.ExtraArgs[DefaultEtcdElectionTimeoutName] = DefaultEtcdElectionTimeoutValue
	}
	if _, ok := c.Core.Etcd.ExtraArgs[DefaultEtcdHeartbeatIntervalName]; !ok {
		c.Core.Etcd.ExtraArgs[DefaultEtcdHeartbeatIntervalName] = DefaultEtcdHeartbeatIntervalValue
	}

	if c.Core.Etcd.BackupConfig != nil {
		if c.Core.Etcd.BackupConfig.IntervalHours == 0 {
			c.Core.Etcd.BackupConfig.IntervalHours = DefaultEtcdBackupConfigIntervalHours
		}
		if c.Core.Etcd.BackupConfig.Retention == 0 {
			c.Core.Etcd.BackupConfig.Retention = DefaultEtcdBackupConfigRetention
		}
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
