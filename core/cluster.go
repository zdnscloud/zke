package core

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/zdnscloud/zke/core/authz"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/core/services"
	"github.com/zdnscloud/zke/pkg/docker"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/k8s"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/util"
	"github.com/zdnscloud/zke/types"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/cert"
)

type Cluster struct {
	AuthnStrategies                    map[string]bool
	ConfigPath                         string
	ConfigDir                          string
	CloudConfigFile                    string
	ControlPlaneHosts                  []*hosts.Host
	Certificates                       map[string]pki.CertificatePKI
	CertificateDir                     string
	ClusterDomain                      string
	ClusterCIDR                        string
	ClusterDNSServer                   string
	DockerDialerFactory                hosts.DialerFactory
	EtcdHosts                          []*hosts.Host
	EtcdReadyHosts                     []*hosts.Host
	InactiveHosts                      []*hosts.Host
	K8sWrapTransport                   k8s.WrapTransport
	KubeClient                         *kubernetes.Clientset
	KubernetesServiceIP                net.IP
	LocalKubeConfigPath                string
	LocalConnDialerFactory             hosts.DialerFactory
	PrivateRegistriesMap               map[string]types.PrivateRegistry
	StateFilePath                      string
	UpdateWorkersOnly                  bool
	UseKubectlDeploy                   bool
	types.ZcloudKubernetesEngineConfig `yaml:",inline"`
	WorkerHosts                        []*hosts.Host
	StorageHosts                       []*hosts.Host
	EdgeHosts                          []*hosts.Host
}

const (
	AuthnX509Provider       = "x509"
	AuthnWebhookProvider    = "webhook"
	StateConfigMapName      = "cluster-state"
	FullStateConfigMapName  = "full-cluster-state"
	UpdateStateTimeout      = 30
	GetStateTimeout         = 30
	KubernetesClientTimeOut = 30
	SyncWorkers             = 10
	NoneAuthorizationMode   = "none"
	LocalNodeAddress        = "127.0.0.1"
	LocalNodeHostname       = "localhost"
	LocalNodeUser           = "root"
	CloudProvider           = "CloudProvider"
	ControlPlane            = "controlPlane"
	WorkerPlane             = "workerPlan"
	EtcdPlane               = "etcd"

	KubeAppLabel = "k8s-app"
	AppLabel     = "app"
	NameLabel    = "name"

	WorkerThreads = util.WorkerThreads
)

func (c *Cluster) DeployControlPlane(ctx context.Context) error {
	// Deploy Etcd Plane
	etcdNodePlanMap := make(map[string]types.ZKEConfigNodePlan)
	// Build etcd node plan map
	for _, etcdHost := range c.EtcdHosts {
		etcdNodePlanMap[etcdHost.Address] = BuildZKEConfigNodePlan(ctx, c, etcdHost, etcdHost.DockerInfo)
	}
	if len(c.Services.Etcd.ExternalURLs) > 0 {
		log.Infof(ctx, "[etcd] External etcd connection string has been specified, skipping etcd plane")
	} else {
		if err := services.RunEtcdPlane(ctx, c.EtcdHosts, etcdNodePlanMap, c.LocalConnDialerFactory, c.PrivateRegistriesMap, c.UpdateWorkersOnly, c.SystemImages.Alpine, c.Services.Etcd, c.Certificates); err != nil {
			return fmt.Errorf("[etcd] Failed to bring up Etcd Plane: %v", err)
		}
	}
	// Deploy Control plane
	cpNodePlanMap := make(map[string]types.ZKEConfigNodePlan)
	// Build cp node plan map
	for _, cpHost := range c.ControlPlaneHosts {
		cpNodePlanMap[cpHost.Address] = BuildZKEConfigNodePlan(ctx, c, cpHost, cpHost.DockerInfo)
	}
	if err := services.RunControlPlane(ctx, c.ControlPlaneHosts,
		c.LocalConnDialerFactory,
		c.PrivateRegistriesMap,
		cpNodePlanMap,
		c.UpdateWorkersOnly,
		c.SystemImages.Alpine,
		c.Certificates); err != nil {
		return fmt.Errorf("[controlPlane] Failed to bring up Control Plane: %v", err)
	}
	return nil
}

func (c *Cluster) DeployWorkerPlane(ctx context.Context) error {
	// Deploy Worker plane
	workerNodePlanMap := make(map[string]types.ZKEConfigNodePlan)
	// Build cp node plan map
	allHosts := hosts.GetUniqueHostList(c.EtcdHosts, c.ControlPlaneHosts, c.WorkerHosts, c.StorageHosts, c.EdgeHosts)
	for _, workerHost := range allHosts {
		workerNodePlanMap[workerHost.Address] = BuildZKEConfigNodePlan(ctx, c, workerHost, workerHost.DockerInfo)
	}
	if err := services.RunWorkerPlane(ctx, allHosts,
		c.LocalConnDialerFactory,
		c.PrivateRegistriesMap,
		workerNodePlanMap,
		c.Certificates,
		c.UpdateWorkersOnly,
		c.SystemImages.Alpine); err != nil {
		return fmt.Errorf("[workerPlane] Failed to bring up Worker Plane: %v", err)
	}
	return nil
}

func ParseConfig(clusterFile string) (*types.ZcloudKubernetesEngineConfig, error) {
	logrus.Debugf("Parsing cluster file [%v]", clusterFile)
	var zkeConfig types.ZcloudKubernetesEngineConfig
	if err := yaml.Unmarshal([]byte(clusterFile), &zkeConfig); err != nil {
		return nil, err
	}
	return &zkeConfig, nil
}

func InitClusterObject(ctx context.Context, zkeConfig *types.ZcloudKubernetesEngineConfig, flags ExternalFlags) (*Cluster, error) {
	// basic cluster object from zkeConfig
	c := &Cluster{
		AuthnStrategies:              make(map[string]bool),
		ZcloudKubernetesEngineConfig: *zkeConfig,
		ConfigPath:                   flags.ClusterFilePath,
		ConfigDir:                    flags.ConfigDir,
		CertificateDir:               flags.CertificateDir,
		StateFilePath:                GetStateFilePath(flags.ClusterFilePath, flags.ConfigDir),
		PrivateRegistriesMap:         make(map[string]types.PrivateRegistry),
	}
	if len(c.ConfigPath) == 0 {
		c.ConfigPath = pki.ClusterConfig
	}
	// set kube_config, state file, and certificate dir
	c.LocalKubeConfigPath = pki.GetLocalKubeConfig(c.ConfigPath, c.ConfigDir)
	c.StateFilePath = GetStateFilePath(c.ConfigPath, c.ConfigDir)
	if len(c.CertificateDir) == 0 {
		c.CertificateDir = GetCertificateDirPath(c.ConfigPath, c.ConfigDir)
	}

	// Setting cluster Defaults
	err := c.setClusterDefaults(ctx)
	if err != nil {
		return nil, err
	}
	// extract cluster network configuration
	c.setNetworkOptions()

	// set hosts groups
	if err := c.InvertIndexHosts(); err != nil {
		return nil, fmt.Errorf("Failed to classify hosts from config file: %v", err)
	}
	// validate cluster configuration
	if err := c.ValidateCluster(); err != nil {
		return nil, fmt.Errorf("Failed to validate cluster: %v", err)
	}
	return c, nil
}

func (c *Cluster) setNetworkOptions() error {
	var err error
	c.KubernetesServiceIP, err = pki.GetKubernetesServiceIP(c.Services.KubeAPI.ServiceClusterIPRange)
	if err != nil {
		return fmt.Errorf("Failed to get Kubernetes Service IP: %v", err)
	}
	c.ClusterDomain = c.Services.Kubelet.ClusterDomain
	c.ClusterCIDR = c.Services.KubeController.ClusterCIDR
	c.ClusterDNSServer = c.Services.Kubelet.ClusterDNSServer
	return nil
}

func (c *Cluster) SetupDialers(ctx context.Context, dailersOptions hosts.DialersOptions) error {
	c.DockerDialerFactory = dailersOptions.DockerDialerFactory
	c.LocalConnDialerFactory = dailersOptions.LocalConnDialerFactory
	c.K8sWrapTransport = dailersOptions.K8sWrapTransport
	// Create k8s wrap transport for bastion host
	if len(c.BastionHost.Address) > 0 {
		var err error
		c.K8sWrapTransport, err = hosts.BastionHostWrapTransport(c.BastionHost)
		if err != nil {
			return err
		}
	}
	return nil
}

func RebuildKubeconfig(ctx context.Context, kubeCluster *Cluster) error {
	return rebuildLocalAdminConfig(ctx, kubeCluster)
}

func rebuildLocalAdminConfig(ctx context.Context, kubeCluster *Cluster) error {
	if len(kubeCluster.ControlPlaneHosts) == 0 {
		return nil
	}
	log.Infof(ctx, "[reconcile] Rebuilding and updating local kube config")
	var workingConfig, newConfig string
	currentKubeConfig := kubeCluster.Certificates[pki.KubeAdminCertName]
	caCrt := kubeCluster.Certificates[pki.CACertName].Certificate
	for _, cpHost := range kubeCluster.ControlPlaneHosts {
		if (currentKubeConfig == pki.CertificatePKI{}) {
			kubeCluster.Certificates = make(map[string]pki.CertificatePKI)
			newConfig = getLocalAdminConfigWithNewAddress(kubeCluster.LocalKubeConfigPath, cpHost.Address, kubeCluster.ClusterName)
		} else {
			kubeURL := fmt.Sprintf("https://%s:6443", cpHost.Address)
			caData := string(cert.EncodeCertPEM(caCrt))
			crtData := string(cert.EncodeCertPEM(currentKubeConfig.Certificate))
			keyData := string(cert.EncodePrivateKeyPEM(currentKubeConfig.Key))
			newConfig = pki.GetKubeConfigX509WithData(kubeURL, kubeCluster.ClusterName, pki.KubeAdminCertName, caData, crtData, keyData)
		}
		if err := pki.DeployAdminConfig(ctx, newConfig, kubeCluster.LocalKubeConfigPath); err != nil {
			return fmt.Errorf("Failed to redeploy local admin config with new host")
		}
		workingConfig = newConfig
		if _, err := GetK8sVersion(kubeCluster.LocalKubeConfigPath, kubeCluster.K8sWrapTransport); err == nil {
			log.Infof(ctx, "[reconcile] host [%s] is active master on the cluster", cpHost.Address)
			break
		}
	}
	currentKubeConfig.Config = workingConfig
	kubeCluster.Certificates[pki.KubeAdminCertName] = currentKubeConfig
	return nil
}

func isLocalConfigWorking(ctx context.Context, localKubeConfigPath string, k8sWrapTransport k8s.WrapTransport) bool {
	if _, err := GetK8sVersion(localKubeConfigPath, k8sWrapTransport); err != nil {
		log.Infof(ctx, "[reconcile] Local config is not valid, rebuilding admin config")
		return false
	}
	return true
}

func getLocalConfigAddress(localConfigPath string) (string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", localConfigPath)
	if err != nil {
		return "", err
	}
	splittedAdress := strings.Split(config.Host, ":")
	address := splittedAdress[1]
	return address[2:], nil
}

func getLocalAdminConfigWithNewAddress(localConfigPath, cpAddress string, clusterName string) string {
	config, _ := clientcmd.BuildConfigFromFlags("", localConfigPath)
	if config == nil {
		return ""
	}
	config.Host = fmt.Sprintf("https://%s:6443", cpAddress)
	return pki.GetKubeConfigX509WithData(
		"https://"+cpAddress+":6443",
		clusterName,
		pki.KubeAdminCertName,
		string(config.CAData),
		string(config.CertData),
		string(config.KeyData))
}

func ApplyAuthzResources(ctx context.Context, zkeConfig types.ZcloudKubernetesEngineConfig, flags ExternalFlags, dailersOptions hosts.DialersOptions) error {
	// dialer factories are not needed here since we are not uses docker only k8s jobs
	kubeCluster, err := InitClusterObject(ctx, &zkeConfig, flags)
	if err != nil {
		return err
	}
	if err := kubeCluster.SetupDialers(ctx, dailersOptions); err != nil {
		return err
	}
	if len(kubeCluster.ControlPlaneHosts) == 0 {
		return nil
	}
	if err := authz.ApplyJobDeployerServiceAccount(ctx, kubeCluster.LocalKubeConfigPath, kubeCluster.K8sWrapTransport); err != nil {
		return fmt.Errorf("Failed to apply the ServiceAccount needed for job execution: %v", err)
	}
	if kubeCluster.Authorization.Mode == NoneAuthorizationMode {
		return nil
	}
	if kubeCluster.Authorization.Mode == services.RBACAuthorizationMode {
		if err := authz.ApplySystemNodeClusterRoleBinding(ctx, kubeCluster.LocalKubeConfigPath, kubeCluster.K8sWrapTransport); err != nil {
			return fmt.Errorf("Failed to apply the ClusterRoleBinding needed for node authorization: %v", err)
		}
	}
	if kubeCluster.Authorization.Mode == services.RBACAuthorizationMode && kubeCluster.Services.KubeAPI.PodSecurityPolicy {
		if err := authz.ApplyDefaultPodSecurityPolicy(ctx, kubeCluster.LocalKubeConfigPath, kubeCluster.K8sWrapTransport); err != nil {
			return fmt.Errorf("Failed to apply default PodSecurityPolicy: %v", err)
		}
		if err := authz.ApplyDefaultPodSecurityPolicyRole(ctx, kubeCluster.LocalKubeConfigPath, kubeCluster.K8sWrapTransport); err != nil {
			return fmt.Errorf("Failed to apply default PodSecurityPolicy ClusterRole and ClusterRoleBinding: %v", err)
		}
	}
	return nil
}

func (c *Cluster) SyncLabelsAndTaints(ctx context.Context, currentCluster *Cluster) error {
	if currentCluster != nil {
		cpToDelete := hosts.GetToDeleteHosts(currentCluster.ControlPlaneHosts, c.ControlPlaneHosts, c.InactiveHosts)
		if len(cpToDelete) == len(currentCluster.ControlPlaneHosts) {
			log.Infof(ctx, "[sync] Cleaning left control plane nodes from reconcilation")
			for _, toDeleteHost := range cpToDelete {
				if err := cleanControlNode(ctx, c, currentCluster, toDeleteHost); err != nil {
					return err
				}
			}
		}
	}

	if len(c.ControlPlaneHosts) > 0 {
		log.Infof(ctx, "[sync] Syncing nodes Labels and Taints")
		k8sClient, err := k8s.NewClient(c.LocalKubeConfigPath, c.K8sWrapTransport)
		if err != nil {
			return fmt.Errorf("Failed to initialize new kubernetes client: %v", err)
		}
		hostList := hosts.GetUniqueHostList(c.EtcdHosts, c.ControlPlaneHosts, c.WorkerHosts, c.StorageHosts, c.EdgeHosts)
		var errgrp errgroup.Group
		hostQueue := make(chan *hosts.Host, len(hostList))
		for _, host := range hostList {
			hostQueue <- host
		}
		close(hostQueue)

		for i := 0; i < SyncWorkers; i++ {
			w := i
			errgrp.Go(func() error {
				var errs []error
				for host := range hostQueue {
					logrus.Debugf("worker [%d] starting sync for node [%s]", w, host.HostnameOverride)
					if err := setNodeAnnotationsLabelsTaints(k8sClient, host); err != nil {
						errs = append(errs, err)
					}
				}
				if len(errs) > 0 {
					return fmt.Errorf("%v", errs)
				}
				return nil
			})
		}
		if err := errgrp.Wait(); err != nil {
			return err
		}
		log.Infof(ctx, "[sync] Successfully synced nodes Labels and Taints")
	}
	return nil
}

func setNodeAnnotationsLabelsTaints(k8sClient *kubernetes.Clientset, host *hosts.Host) error {
	node := &v1.Node{}
	var err error
	for retries := 0; retries <= 5; retries++ {
		node, err = k8s.GetNode(k8sClient, host.HostnameOverride)
		if err != nil {
			logrus.Debugf("[hosts] Can't find node by name [%s], retrying..", host.HostnameOverride)
			time.Sleep(2 * time.Second)
			continue
		}

		oldNode := node.DeepCopy()
		k8s.SetNodeAddressesAnnotations(node, host.InternalAddress, host.Address)
		k8s.SyncNodeLabels(node, host.ToAddLabels, host.ToDelLabels)
		k8s.SyncNodeTaints(node, host.ToAddTaints, host.ToDelTaints)

		if reflect.DeepEqual(oldNode, node) {
			logrus.Debugf("skipping syncing labels for node [%s]", node.Name)
			return nil
		}
		_, err = k8sClient.CoreV1().Nodes().Update(node)
		if err != nil {
			logrus.Debugf("Error syncing labels for node [%s]: %v", node.Name, err)
			time.Sleep(5 * time.Second)
			continue
		}
		return nil
	}
	return err
}

func (c *Cluster) PrePullK8sImages(ctx context.Context) error {
	log.Infof(ctx, "Pre-pulling kubernetes images")
	var errgrp errgroup.Group
	hostList := hosts.GetUniqueHostList(c.EtcdHosts, c.ControlPlaneHosts, c.WorkerHosts, c.StorageHosts, c.EdgeHosts)
	hostsQueue := util.GetObjectQueue(hostList)
	for w := 0; w < WorkerThreads; w++ {
		errgrp.Go(func() error {
			var errList []error
			for host := range hostsQueue {
				runHost := host.(*hosts.Host)
				err := docker.UseLocalOrPull(ctx, runHost.DClient, runHost.Address, c.SystemImages.Kubernetes, "pre-deploy", c.PrivateRegistriesMap)
				if err != nil {
					errList = append(errList, err)
				}
			}
			return util.ErrList(errList)
		})
	}

	if err := errgrp.Wait(); err != nil {
		return err
	}
	log.Infof(ctx, "Kubernetes images pulled successfully")
	return nil
}

func RestartClusterPods(ctx context.Context, kubeCluster *Cluster) error {
	log.Infof(ctx, "Restarting network, ingress, and metrics pods")
	// this will remove the pods created by ZKE and let the controller creates them again
	kubeClient, err := k8s.NewClient(kubeCluster.LocalKubeConfigPath, kubeCluster.K8sWrapTransport)
	if err != nil {
		return fmt.Errorf("Failed to initialize new kubernetes client: %v", err)
	}
	labelsList := []string{
		fmt.Sprintf("%s=%s", KubeAppLabel, CalicoNetworkPlugin),
		fmt.Sprintf("%s=%s", KubeAppLabel, FlannelNetworkPlugin),
		fmt.Sprintf("%s=%s", AppLabel, NginxIngressAddonAppName),
		fmt.Sprintf("%s=%s", KubeAppLabel, DefaultMonitoringMetricsProvider),
		fmt.Sprintf("%s=%s", KubeAppLabel, CoreDNSAddonAppName),
	}
	var errgrp errgroup.Group
	labelQueue := util.GetObjectQueue(labelsList)
	for w := 0; w < services.WorkerThreads; w++ {
		errgrp.Go(func() error {
			var errList []error
			for label := range labelQueue {
				runLabel := label.(string)
				// list pods to be deleted
				pods, err := k8s.ListPodsByLabel(kubeClient, runLabel)
				if err != nil {
					errList = append(errList, err)
				}
				// delete pods
				err = k8s.DeletePods(kubeClient, pods)
				if err != nil {
					errList = append(errList, err)
				}
			}
			return util.ErrList(errList)
		})
	}
	if err := errgrp.Wait(); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) GetHostInfoMap() map[string]dockertypes.Info {
	hostsInfoMap := make(map[string]dockertypes.Info)
	allHosts := hosts.GetUniqueHostList(c.EtcdHosts, c.ControlPlaneHosts, c.WorkerHosts, c.StorageHosts, c.EdgeHosts)
	for _, host := range allHosts {
		hostsInfoMap[host.Address] = host.DockerInfo
	}
	return hostsInfoMap
}