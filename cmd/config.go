package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/core/services"
	"github.com/zdnscloud/zke/types"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

const (
	comments = `# If you intened to deploy Kubernetes in an air-gapped environment,
# please consult the documentation on how to configure custom ZKE images.`
	DefaultClusterSSHKeyPath     = "~/.ssh/id_rsa"
	DefaultClusterSSHKey         = ""
	DefaultClusterSSHPort        = "22"
	DefaultClusterSSHUser        = "ubuntu"
	DefaultClusterDockerSockPath = "/var/run/docker.sock"

	IngressSelectLabel = "node-role.kubernetes.io/edge"
)

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
		},
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

func writeConfig(cluster *types.ZKEConfig, configFile string) error {
	yamlConfig, err := yaml.Marshal(*cluster)
	if err != nil {
		return err
	}
	logrus.Debugf("Deploying cluster configuration file: %s", configFile)
	configString := fmt.Sprintf("%s\n%s", comments, string(yamlConfig))
	return ioutil.WriteFile(configFile, []byte(configString), 0640)
}

func clusterConfig(ctx *cli.Context) error {
	configFile := ctx.String("name")
	cluster := types.ZKEConfig{}
	// set zke config version
	cluster.Version = defaultConfigVersion
	// Get cluster config from user
	reader := bufio.NewReader(os.Stdin)
	// Generate empty configuration file
	if ctx.Bool("empty") {
		cluster.Nodes = make([]types.ZKEConfigNode, 1)
		return writeConfig(&cluster, configFile)
	}
	sshKeyPath, err := getConfig(reader, "Cluster Level SSH Private Key Path", DefaultClusterSSHKeyPath)
	if err != nil {
		return err
	}
	cluster.Option.SSHKeyPath = sshKeyPath
	sshPort, err := getConfig(reader, "Cluster Level SSH Port of all host", DefaultClusterSSHPort)
	if err != nil {
		return err
	}
	cluster.Option.SSHPort = sshPort
	sshUser, err := getConfig(reader, "Cluster Level SSH User of all host", DefaultClusterSSHUser)
	if err != nil {
		return err
	}
	cluster.Option.SSHUser = sshUser
	dockerSocketPath, err := getConfig(reader, "Cluster Level Docker socket path on all host", DefaultClusterDockerSockPath)
	if err != nil {
		return err
	}
	cluster.Option.DockerSocket = dockerSocketPath

	clusterDomain, err := getConfig(reader, "Cluster domain", core.DefaultClusterDomain)
	if err != nil {
		return err
	}
	cluster.Option.ClusterDomain = clusterDomain

	serviceClusterIPRange, err := getConfig(reader, "Service Cluster IP Range", core.DefaultServiceClusterIPRange)
	if err != nil {
		return err
	}
	cluster.Option.ServiceClusterIpRange = serviceClusterIPRange

	clusterNetworkCidr, err := getConfig(reader, "Cluster Network CIDR", core.DefaultClusterCIDR)
	if err != nil {
		return err
	}
	cluster.Option.ClusterCidr = clusterNetworkCidr

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
		hostCfg, err := getHostConfig(reader, i, &cluster)
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
	cluster.Network.Ingress = *ingressConfig

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
	cluster.Image = types.K8sVersionToZKESystemImages[core.DefaultK8sVersion]
	cluster.Network.DNS.UpstreamNameservers, err = getGlobalDNSConfig(reader)
	if err != nil {
		return err
	}
	// Get Services Config
	coreServiceConfig, err := getServiceConfig(reader, &cluster)
	if err != nil {
		return err
	}
	cluster.Core = *coreServiceConfig

	return writeConfig(&cluster, configFile)
}

func getHostConfig(reader *bufio.Reader, index int, cluster *types.ZKEConfig) (*types.ZKEConfigNode, error) {
	host := types.ZKEConfigNode{}
	address, err := getConfig(reader, fmt.Sprintf("SSH Address of host (%d)", index+1), "")
	if err != nil {
		return nil, err
	}
	host.Address = address
	host.Port = cluster.Option.SSHPort
	host.User = cluster.Option.SSHUser
	host.SSHKey = cluster.Option.SSHKey
	host.SSHKeyPath = cluster.Option.SSHKeyPath
	host.DockerSocket = cluster.Option.DockerSocket

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

	isEdgeHost, err := getConfig(reader, fmt.Sprintf("Is host (%s) an Edge host (y/n)?", address), "y")
	if err != nil {
		return nil, err
	}
	if isEdgeHost == "y" || isEdgeHost == "Y" {
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

func getServiceConfig(reader *bufio.Reader, cluster *types.ZKEConfig) (*types.ZKEConfigCore, error) {
	servicesConfig := types.ZKEConfigCore{}
	servicesConfig.Etcd = types.ETCDService{}
	servicesConfig.KubeAPI = types.KubeAPIService{}
	servicesConfig.KubeController = types.KubeControllerService{}
	servicesConfig.Scheduler = types.SchedulerService{}
	servicesConfig.Kubelet = types.KubeletService{}
	servicesConfig.Kubeproxy = types.KubeproxyService{}
	servicesConfig.Kubelet.ClusterDomain = cluster.Option.ClusterDomain
	servicesConfig.KubeAPI.ServiceClusterIPRange = cluster.Option.ServiceClusterIpRange
	servicesConfig.KubeController.ServiceClusterIPRange = cluster.Option.ServiceClusterIpRange
	podSecurityPolicy, err := getConfig(reader, "Enable PodSecurityPolicy", "n")
	if err != nil {
		return nil, err
	}
	if podSecurityPolicy == "y" || podSecurityPolicy == "Y" {
		servicesConfig.KubeAPI.PodSecurityPolicy = true
	} else {
		servicesConfig.KubeAPI.PodSecurityPolicy = false
	}
	servicesConfig.KubeController.ClusterCIDR = cluster.Option.ClusterCidr
	clusterDNSServiceIP, err := getConfig(reader, "Cluster DNS Service IP", core.DefaultClusterDNSService)
	if err != nil {
		return nil, err
	}
	servicesConfig.Kubelet.ClusterDNSServer = clusterDNSServiceIP
	return &servicesConfig, nil
}

func getAuthnConfig(reader *bufio.Reader) (*types.ZKEConfigAuthn, error) {
	authnConfig := types.ZKEConfigAuthn{}
	authnType, err := getConfig(reader, "Authentication Strategy", core.DefaultAuthStrategy)
	if err != nil {
		return nil, err
	}
	authnConfig.Strategy = authnType
	return &authnConfig, nil
}

func getAuthzConfig(reader *bufio.Reader) (*types.ZKEConfigAuthz, error) {
	authzConfig := types.ZKEConfigAuthz{}
	authzMode, err := getConfig(reader, "Authorization Mode (rbac, none)", core.DefaultAuthorizationMode)
	if err != nil {
		return nil, err
	}
	authzConfig.Mode = authzMode
	return &authzConfig, nil
}

func getNetworkConfig(reader *bufio.Reader) (*types.ZKEConfigNetwork, error) {
	networkConfig := types.ZKEConfigNetwork{}
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
		networkConfig.Iface = networkFlannelIface
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
