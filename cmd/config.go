package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/core/services"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/util"
	"github.com/zdnscloud/zke/types"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

const (
	comments = `# If you intened to deploy Kubernetes in an air-gapped environment,
# please consult the documentation on how to configure custom ZKE images.`
	FlannelIface                 = "flannel_iface"
	FlannelBackendType           = "flannel_backend_type"
	FlannelBackendDirectrouting  = "flannel_vxlan_directrouting"
	DefaultClusterSSHKeyPath     = "~/.ssh/id_rsa"
	DefaultClusterSSHKey         = ""
	DefaultClusterSSHPort        = "22"
	DefaultClusterSSHUser        = "ubuntu"
	DefaultClusterDockerSockPath = "/var/run/docker.sock"

	IngressSelectLabel = "node-role.kubernetes.io/edge"
)

type clusterCommonCfg struct {
	sshPort      string
	sshKeyPath   string
	sshUser      string
	dockerSocket string
}

func ConfigCommand() cli.Command {
	return cli.Command{
		Name:   "config",
		Usage:  "Setup cluster configuration",
		Action: clusterConfig,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name,n",
				Usage: "Name of the configuration file",
				Value: pki.ClusterConfig,
			},
			cli.BoolFlag{
				Name:  "empty,e",
				Usage: "Generate Empty configuration file",
			},
			cli.BoolFlag{
				Name:  "print,p",
				Usage: "Print configuration",
			},
			cli.BoolFlag{
				Name:  "system-images",
				Usage: "Generate the default system images",
			},
			cli.BoolFlag{
				Name:  "all",
				Usage: "Generate the default system images for all versions",
			},
			cli.StringFlag{
				Name:  "version",
				Usage: "Generate the default system images for specific k8s versions",
			},
		},
	}
}

func GetConfig(reader *bufio.Reader, text, def string) (string, error) {
	for {
		if def == "" {
			fmt.Printf("[+] %s [%s]: ", text, "none")
		} else {
			fmt.Printf("[+] %s [%s]: ", text, def)
		}
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)
		if input != "" {
			return input, nil
		}
		return def, nil
	}
}
func getConfig(reader *bufio.Reader, text, def string) (string, error) {
	for {
		if def == "" {
			fmt.Printf("[+] %s [%s]: ", text, "none")
		} else {
			fmt.Printf("[+] %s [%s]: ", text, def)
		}
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)
		if input != "" {
			return input, nil
		}
		return def, nil
	}
}

func writeConfig(cluster *types.ZcloudKubernetesEngineConfig, configFile string, print bool) error {
	yamlConfig, err := yaml.Marshal(*cluster)
	if err != nil {
		return err
	}
	logrus.Debugf("Deploying cluster configuration file: %s", configFile)
	configString := fmt.Sprintf("%s\n%s", comments, string(yamlConfig))
	if print {
		fmt.Printf("Configuration File: \n%s", configString)
		return nil
	}
	return ioutil.WriteFile(configFile, []byte(configString), 0640)
}

func clusterConfig(ctx *cli.Context) error {
	if ctx.Bool("system-images") {
		return generateSystemImagesList(ctx.String("version"), ctx.Bool("all"))
	}
	configFile := ctx.String("name")
	print := ctx.Bool("print")
	cluster := types.ZcloudKubernetesEngineConfig{}
	// set zke config version
	cluster.ConfigVersion = defaultConfigVersion
	// Get cluster config from user
	reader := bufio.NewReader(os.Stdin)
	// Generate empty configuration file
	if ctx.Bool("empty") {
		cluster.Nodes = make([]types.ZKEConfigNode, 1)
		return writeConfig(&cluster, configFile, print)
	}
	sshKeyPath, err := getConfig(reader, "Cluster Level SSH Private Key Path", DefaultClusterSSHKeyPath)
	if err != nil {
		return err
	}
	cluster.SSHKeyPath = sshKeyPath
	sshPort, err := getConfig(reader, "Cluster Level SSH Port of all host", DefaultClusterSSHPort)
	if err != nil {
		return err
	}
	cluster.SSHPort = sshPort
	sshUser, err := getConfig(reader, "Cluster Level SSH User of all host", DefaultClusterSSHUser)
	if err != nil {
		return err
	}
	cluster.SSHUser = sshUser
	dockerSocketPath, err := getConfig(reader, "Cluster Level Docker socket path on all host", DefaultClusterDockerSockPath)
	if err != nil {
		return err
	}
	cluster.DockerSocket = dockerSocketPath
	hostCommonCfg := clusterCommonCfg{sshPort, sshKeyPath, sshUser, dockerSocketPath}
	// Get number of hosts
	numberOfHostsString, err := getConfig(reader, "Number of Hosts", "1")
	if err != nil {
		return err
	}
	numberOfHostsInt, err := strconv.Atoi(numberOfHostsString)
	if err != nil {
		return err
	}
	// Get Hosts config
	cluster.Nodes = make([]types.ZKEConfigNode, 0)
	for i := 0; i < numberOfHostsInt; i++ {
		hostCfg, err := getHostConfig(reader, i, hostCommonCfg)
		if err != nil {
			return err
		}
		cluster.Nodes = append(cluster.Nodes, *hostCfg)
	}
	// Get Network config
	networkConfig, err := getNetworkConfig(reader)
	if err != nil {
		return err
	}
	cluster.Network = *networkConfig
	// Get Ingress config
	ingressConfig, err := getIngressConfig(reader, cluster.Nodes)
	if err != nil {
		return err
	}
	cluster.Ingress = *ingressConfig
	// Get Storage config
	storageConfig, err := getStorageConfig(reader, cluster.Nodes)
	if err != nil {
		return err
	}
	cluster.Storage = *storageConfig

	// Get Authentication Config
	authnConfig, err := getAuthnConfig(reader)
	if err != nil {
		return err
	}
	cluster.Authentication = *authnConfig
	// Get Authorization config
	authzConfig, err := getAuthzConfig(reader)
	if err != nil {
		return err
	}
	cluster.Authorization = *authzConfig
	// Get k8s/system images
	systemImages, err := getSystemImagesConfig(reader)
	if err != nil {
		return err
	}
	cluster.SystemImages = *systemImages
	cluster.DNS.UpstreamNameservers, err = getGlobalDNSConfig(reader)
	if err != nil {
		return err
	}
	// Get Services Config
	serviceConfig, err := getServiceConfig(reader)
	if err != nil {
		return err
	}
	cluster.Services = *serviceConfig
	cluster.Monitoring.GrafanaIngressEndpoint = "grafana.kube-monitoring." + cluster.Services.Kubelet.ClusterDomain
	cluster.Monitoring.PrometheusAlertManagerIngressEndpoint = "alertmanager.kube-monitoring." + cluster.Services.Kubelet.ClusterDomain
	registryConfig, err := getRegistryConfig(reader, &cluster)
	if err != nil {
		return err
	}
	cluster.Registry = *registryConfig
	return writeConfig(&cluster, configFile, print)
}

func getHostConfig(reader *bufio.Reader, index int, hostCommonCfg clusterCommonCfg) (*types.ZKEConfigNode, error) {
	host := types.ZKEConfigNode{}
	address, err := getConfig(reader, fmt.Sprintf("SSH Address of host (%d)", index+1), "")
	if err != nil {
		return nil, err
	}
	host.Address = address
	host.Port = hostCommonCfg.sshPort
	host.User = hostCommonCfg.sshUser
	host.SSHKey = DefaultClusterSSHKey
	host.SSHKeyPath = hostCommonCfg.sshKeyPath
	host.DockerSocket = hostCommonCfg.dockerSocket
	isControlHost, err := getConfig(reader, fmt.Sprintf("Is host (%s) a Control Plane host (y/n)?", address), "y")
	if err != nil {
		return nil, err
	}
	if isControlHost == "y" || isControlHost == "Y" {
		host.Role = append(host.Role, services.ControlRole)
	}
	isWorkerHost, err := getConfig(reader, fmt.Sprintf("Is host (%s) a Worker host (y/n)?", address), "n")
	if err != nil {
		return nil, err
	}
	if isWorkerHost == "y" || isWorkerHost == "Y" {
		host.Role = append(host.Role, services.WorkerRole)
	}
	isEtcdHost, err := getConfig(reader, fmt.Sprintf("Is host (%s) an etcd host (y/n)?", address), "n")
	if err != nil {
		return nil, err
	}
	if isEtcdHost == "y" || isEtcdHost == "Y" {
		host.Role = append(host.Role, services.ETCDRole)
	}
	isStorageHost, err := getConfig(reader, fmt.Sprintf("Is host (%s) an Storage host (y/n)?", address), "n")
	if err != nil {
		return nil, err
	}
	if isStorageHost == "y" || isStorageHost == "Y" {
		host.Role = append(host.Role, services.StorageRole)
	}
	isNetBorderHost, err := getConfig(reader, fmt.Sprintf("Is host (%s) an Edge host (y/n)?", address), "y")
	if err != nil {
		return nil, err
	}
	if isNetBorderHost == "y" || isNetBorderHost == "Y" {
		host.Role = append(host.Role, services.EdgeRole)
	}
	hostnameOverride, err := getConfig(reader, fmt.Sprintf("Override Hostname of host (%s)", address), "")
	if err != nil {
		return nil, err
	}
	host.HostnameOverride = hostnameOverride
	internalAddress, err := getConfig(reader, fmt.Sprintf("Internal IP of host (%s)", address), "")
	if err != nil {
		return nil, err
	}
	host.InternalAddress = internalAddress
	return &host, nil
}

func getSystemImagesConfig(reader *bufio.Reader) (*types.ZKESystemImages, error) {
	imageDefaults := types.K8sVersionToZKESystemImages[core.DefaultK8sVersion]
	kubeImage, err := getConfig(reader, "Kubernetes Docker image", imageDefaults.Kubernetes)
	if err != nil {
		return nil, err
	}
	systemImages, ok := types.K8sVersionToZKESystemImages[kubeImage]
	if ok {
		return &systemImages, nil
	}
	imageDefaults.Kubernetes = kubeImage
	return &imageDefaults, nil
}

func getServiceConfig(reader *bufio.Reader) (*types.ZKEConfigServices, error) {
	servicesConfig := types.ZKEConfigServices{}
	servicesConfig.Etcd = types.ETCDService{}
	servicesConfig.KubeAPI = types.KubeAPIService{}
	servicesConfig.KubeController = types.KubeControllerService{}
	servicesConfig.Scheduler = types.SchedulerService{}
	servicesConfig.Kubelet = types.KubeletService{}
	servicesConfig.Kubeproxy = types.KubeproxyService{}
	clusterDomain, err := getConfig(reader, "Cluster domain", core.DefaultClusterDomain)
	if err != nil {
		return nil, err
	}
	servicesConfig.Kubelet.ClusterDomain = clusterDomain
	serviceClusterIPRange, err := getConfig(reader, "Service Cluster IP Range", core.DefaultServiceClusterIPRange)
	if err != nil {
		return nil, err
	}
	servicesConfig.KubeAPI.ServiceClusterIPRange = serviceClusterIPRange
	servicesConfig.KubeController.ServiceClusterIPRange = serviceClusterIPRange
	podSecurityPolicy, err := getConfig(reader, "Enable PodSecurityPolicy", "n")
	if err != nil {
		return nil, err
	}
	if podSecurityPolicy == "y" || podSecurityPolicy == "Y" {
		servicesConfig.KubeAPI.PodSecurityPolicy = true
	} else {
		servicesConfig.KubeAPI.PodSecurityPolicy = false
	}
	clusterNetworkCidr, err := getConfig(reader, "Cluster Network CIDR", core.DefaultClusterCIDR)
	if err != nil {
		return nil, err
	}
	servicesConfig.KubeController.ClusterCIDR = clusterNetworkCidr
	clusterDNSServiceIP, err := getConfig(reader, "Cluster DNS Service IP", core.DefaultClusterDNSService)
	if err != nil {
		return nil, err
	}
	servicesConfig.Kubelet.ClusterDNSServer = clusterDNSServiceIP
	return &servicesConfig, nil
}

func getAuthnConfig(reader *bufio.Reader) (*types.AuthnConfig, error) {
	authnConfig := types.AuthnConfig{}
	authnType, err := getConfig(reader, "Authentication Strategy", core.DefaultAuthStrategy)
	if err != nil {
		return nil, err
	}
	authnConfig.Strategy = authnType
	return &authnConfig, nil
}

func getAuthzConfig(reader *bufio.Reader) (*types.AuthzConfig, error) {
	authzConfig := types.AuthzConfig{}
	authzMode, err := getConfig(reader, "Authorization Mode (rbac, none)", core.DefaultAuthorizationMode)
	if err != nil {
		return nil, err
	}
	authzConfig.Mode = authzMode
	return &authzConfig, nil
}

func getNetworkConfig(reader *bufio.Reader) (*types.NetworkConfig, error) {
	networkConfig := types.NetworkConfig{}
	networkPlugin, err := getConfig(reader, "Network Plugin Type (flannel, calico)", core.DefaultNetworkPlugin)
	if err != nil {
		return nil, err
	}
	networkConfig.Plugin = networkPlugin
	if networkPlugin == core.DefaultNetworkPlugin {
		networkFlannelIface, err := getConfig(reader, "Flannel Network Interface", "")
		if err != nil {
			return nil, err
		}
		networkConfig.Options = make(map[string]string)
		networkConfig.Options[FlannelIface] = networkFlannelIface
		networkFlannelBackendType, err := getConfig(reader, "Flannel Backend Type (vxlan, host-gw)", core.DefaultFlannelBackendType)
		if err != nil {
			return nil, err
		}
		networkConfig.Options[FlannelBackendType] = networkFlannelBackendType
		if networkFlannelBackendType == core.DefaultFlannelBackendType {
			networkConfig.Options[FlannelBackendDirectrouting] = "true"
		} else {
			networkConfig.Options[FlannelBackendDirectrouting] = "false"
		}
	}
	return &networkConfig, nil
}

func getIngressConfig(reader *bufio.Reader, nodes []types.ZKEConfigNode) (*types.IngressConfig, error) {
	ingressCfg := types.IngressConfig{}
	ingressCfg.NodeSelector = make(map[string]string)
	for _, n := range nodes {
		for _, v := range n.Role {
			if v == services.EdgeRole {
				ingressCfg.NodeSelector[IngressSelectLabel] = "true"
			}
		}
	}
	return &ingressCfg, nil
}

func getStorageConfig(reader *bufio.Reader, nodes []types.ZKEConfigNode) (*types.StorageConfig, error) {
	storageinfo := make(map[string][]string)
	for _, n := range nodes {
		for _, v := range n.Role {
			if v == services.StorageRole {
				var Host string
				devices, err := getConfig(reader, fmt.Sprintf("Storage disk partitions on host (%s),separated by commas", n.Address), "")
				if err != nil {
					return nil, err
				}
				if n.HostnameOverride != "" {
					Host = n.HostnameOverride
				} else {
					Host = n.Address
				}
				storageinfo[Host] = strings.Split(devices, ",")
			}
		}
	}
	storageCfg := types.StorageConfig{}
	hostsequence := promptStorage(storageinfo)
	storagetypes := []string{"Lvm", "Nfs", "Ceph"}
	for _, t := range storagetypes {
		cfg, err := allocateStorage(reader, storageinfo, t, hostsequence)
		if err != nil {
			return nil, err
		}
		switch t {
		case "Lvm":
			storageCfg.Lvm = cfg
		case "Nfs":
			storageCfg.Nfs = cfg
		case "Ceph":
			storageCfg.Ceph = cfg
		}
	}
	return &storageCfg, nil
}

func promptStorage(storageinfo map[string][]string) map[int]string {
	hostsequence := make(map[int]string)
	var keys []string
	for k := range storageinfo {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		fmt.Printf("%d: %s\n", i+1, k)
		hostsequence[i+1] = k
	}
	return hostsequence
}

func allocateStorage(reader *bufio.Reader, storageinfo map[string][]string, t string, hostsequence map[int]string) ([]types.Deviceconf, error) {
	hosts := make([]string, 0)
	for {
		hostnums, err := GetConfig(reader, fmt.Sprintf("Host number used for %s storage,separated by commas", t), "")
		if err != nil {
			return nil, err
		}
		if hostnums == "" {
			break
		}
		nums := strings.Split(strings.TrimSuffix(hostnums, ","), ",")
		var flag bool
		var validnum []int
		for _, v := range nums {
			num, _ := strconv.Atoi(v)
			_, ok := hostsequence[num]
			if ok {
				validnum = append(validnum, num)
				flag = true
			} else {
				fmt.Printf("The host not exist or has been allocated, please change input!\n")
				validnum = validnum[0:0]
				flag = false
				break
			}
		}
		if flag {
			for _, v := range validnum {
				hosts = append(hosts, hostsequence[v])
				delete(hostsequence, v)
			}
			break
		} else {
			continue
		}
	}
	devicecfgs := make([]types.Deviceconf, 0)
	for _, h := range hosts {
		c := types.Deviceconf{
			Host: h,
			Devs: storageinfo[h],
		}
		devicecfgs = append(devicecfgs, c)
	}
	return devicecfgs, nil
}

func getGlobalDNSConfig(reader *bufio.Reader) ([]string, error) {
	globalDNS := []string{}
	inputString, err := getConfig(reader, fmt.Sprintf("Cluster global dns,separated by commas"), core.DefaultClusterGlobalDns)
	if err != nil {
		return nil, err
	}
	servers := strings.Split(inputString, ",")
	for _, server := range servers {
		globalDNS = append(globalDNS, server)
	}
	return globalDNS, nil
}

func getRegistryConfig(reader *bufio.Reader, c *types.ZcloudKubernetesEngineConfig) (*types.RegistryConfig, error) {
	ctx := context.TODO()
	registryCfg := types.RegistryConfig{}
	isenabled, err := getConfig(reader, fmt.Sprintf("Is enabled harbor registry (y/n)?"), "y")
	if err != nil {
		return nil, err
	}
	if isenabled == "n" || isenabled == "N" {
		registryCfg.Isenabled = false
		return &registryCfg, nil
	}
	if isenabled == "y" || isenabled == "Y" {
		if len(c.Storage.Lvm) == 0 {
			log.Warnf(ctx, "None available lvm storge, will not enable harbor registry!")
			registryCfg.Isenabled = false
			return &registryCfg, nil
		}
		registryCfg.Isenabled = true
		registryDiskCapacity, err := getConfig(reader, fmt.Sprintf("Cluster registry disk capacity"), "50Gi")
		if err != nil {
			return nil, err
		}
		registryCfg.RegistryDiskCapacity = registryDiskCapacity
		registryIngressURL, err := getConfig(reader, fmt.Sprintf("Cluster registry ingress url"), "registry.kube-registry."+c.Services.Kubelet.ClusterDomain)
		if err != nil {
			return nil, err
		}
		registryCfg.RegistryIngressURL = registryIngressURL
		registryCfg.NotaryIngressURL = "notary.kube-registry." + c.Services.Kubelet.ClusterDomain
		registryCfg.RedisDiskCapacity = core.DefaultRegistryRedisDiskCapacity
		registryCfg.DatabaseDiskCapacity = core.DefaultRegistryDatabaseDiskCapacity
		registryCfg.JobserviceDiskCapacity = core.DefaultRegistryJobserviceDiskCapacity
		registryCfg.ChartmuseumDiskCapacity = core.DefaultRegistryChartmuseumDiskCapacity
	}
	return &registryCfg, nil
}

func generateSystemImagesList(version string, all bool) error {
	allVersions := []string{}
	currentVersionImages := make(map[string]types.ZKESystemImages)
	for version := range types.AllK8sVersions {
		err := util.ValidateVersion(version)
		if err != nil {
			continue
		}
		allVersions = append(allVersions, version)
		currentVersionImages[version] = types.AllK8sVersions[version]
	}
	if all {
		for version, zkeSystemImages := range currentVersionImages {
			err := util.ValidateVersion(version)
			if err != nil {
				continue
			}
			logrus.Infof("Generating images list for version [%s]:", version)
			uniqueImages := getUniqueSystemImageList(zkeSystemImages)
			for _, image := range uniqueImages {
				if image == "" {
					continue
				}
				fmt.Printf("%s\n", image)
			}
		}
		return nil
	}
	if len(version) == 0 {
		version = types.DefaultK8s
	}
	zkeSystemImages := types.AllK8sVersions[version]
	if zkeSystemImages == (types.ZKESystemImages{}) {
		return fmt.Errorf("k8s version is not supported, supported versions are: %v", allVersions)
	}
	logrus.Infof("Generating images list for version [%s]:", version)
	uniqueImages := getUniqueSystemImageList(zkeSystemImages)
	for _, image := range uniqueImages {
		if image == "" {
			continue
		}
		fmt.Printf("%s\n", image)
	}
	return nil
}

func getUniqueSystemImageList(zkeSystemImages types.ZKESystemImages) []string {
	imagesReflect := reflect.ValueOf(zkeSystemImages)
	images := make([]string, imagesReflect.NumField())
	for i := 0; i < imagesReflect.NumField(); i++ {
		images[i] = imagesReflect.Field(i).Interface().(string)
	}
	return getUniqueSlice(images)
}

func getUniqueSlice(slice []string) []string {
	encountered := map[string]bool{}
	unqiue := []string{}
	for i := range slice {
		if encountered[slice[i]] {
			continue
		} else {
			encountered[slice[i]] = true
			unqiue = append(unqiue, slice[i])
		}
	}
	return unqiue
}
