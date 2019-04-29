package network

import (
	"context"
	"fmt"

	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/network/calico"
	"github.com/zdnscloud/zke/network/coredns"
	"github.com/zdnscloud/zke/network/flannel"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pki"
	"github.com/zdnscloud/zke/templates"
)

const (
	ClusterVersion = "ClusterVersion"
	ClusterCIDR    = "ClusterCIDR"
	CloudProvider  = "CloudProvider"
	RBACConfig     = "RBACConfig"
	KubeCfg        = "KubeCfg"

	NetworkPluginResourceName = "zke-network-plugin"
	NoNetworkPlugin           = "none"

	FlannelNetworkPlugin        = "flannel"
	FlannelIface                = "flannel_iface"
	FlannelBackendType          = "flannel_backend_type"
	FlannelBackendDirectrouting = "flannel_vxlan_directrouting"
	FlannelInterface            = "FlannelInterface"
	FlannelBackend              = "FlannelBackend"

	CalicoNetworkPlugin = "calico"
	CalicoCloudProvider = "calico_cloud_provider"
	Calicoctl           = "Calicoctl"

	CoreDNSResourceName = "zke-dns-plugin"

	Image            = "Image"
	CNIImage         = "CNIImage"
	NodeImage        = "NodeImage"
	ControllersImage = "ControllersImage"
)

func DeployNetwork(ctx context.Context, c *cluster.Cluster) error {
	if err := DeployNetworkPlugin(ctx, c); err != nil {
		return err
	}

	if err := deployDNSPlugin(ctx, c); err != nil {
		return err
	}
	return nil
}

func DeployNetworkPlugin(ctx context.Context, c *cluster.Cluster) error {
	log.Infof(ctx, "[network] Setting up network plugin: %s", c.Network.Plugin)
	switch c.Network.Plugin {
	case FlannelNetworkPlugin:
		return doFlannelDeploy(ctx, c)
	case CalicoNetworkPlugin:
		return doCalicoDeploy(ctx, c)
	case NoNetworkPlugin:
		log.Infof(ctx, "[network] Not deploying a cluster network, expecting custom CNI")
		return nil
	default:
		return fmt.Errorf("[network] Unsupported network plugin: %s", c.Network.Plugin)
	}
}

func doFlannelDeploy(ctx context.Context, c *cluster.Cluster) error {
	flannelConfig := map[string]interface{}{
		ClusterCIDR:      c.ClusterCIDR,
		Image:            c.SystemImages.Flannel,
		CNIImage:         c.SystemImages.FlannelCNI,
		FlannelInterface: c.Network.Options[FlannelIface],
		FlannelBackend: map[string]interface{}{
			"Type":          c.Network.Options[FlannelBackendType],
			"Directrouting": c.Network.Options[FlannelBackendDirectrouting],
		},
		RBACConfig:     c.Authorization.Mode,
		ClusterVersion: cluster.GetTagMajorVersion(c.Version),
	}
	//pluginYaml, err := templates.GetManifest(flannelConfig, FlannelNetworkPlugin)
	pluginYaml, err := templates.CompileTemplateFromMap(flannel.FlannelTemplate, flannelConfig)
	if err != nil {
		return err
	}
	return c.DoAddonDeploy(ctx, pluginYaml, NetworkPluginResourceName, true)
}

func doCalicoDeploy(ctx context.Context, c *cluster.Cluster) error {
	clientConfig := pki.GetConfigPath(pki.KubeNodeCertName)
	calicoConfig := map[string]interface{}{
		KubeCfg:       clientConfig,
		ClusterCIDR:   c.ClusterCIDR,
		CNIImage:      c.SystemImages.CalicoCNI,
		NodeImage:     c.SystemImages.CalicoNode,
		Calicoctl:     c.SystemImages.CalicoCtl,
		CloudProvider: c.Network.Options[CalicoCloudProvider],
		RBACConfig:    c.Authorization.Mode,
	}

	var CalicoTemplate string
	switch c.Version {
	case "v1.13.1":
		CalicoTemplate = calico.CalicoTemplateV113
	case "default":
		CalicoTemplate = calico.CalicoTemplateV112
	}
	//pluginYaml, err := templates.GetManifest(calicoConfig, CalicoNetworkPlugin, c.Version)
	pluginYaml, err := templates.CompileTemplateFromMap(CalicoTemplate, calicoConfig)
	if err != nil {
		return err
	}
	return c.DoAddonDeploy(ctx, pluginYaml, NetworkPluginResourceName, true)
}

func deployDNSPlugin(ctx context.Context, c *cluster.Cluster) error {
	if err := doCoreDNSDeploy(ctx, c); err != nil {
		if err, ok := err.(*cluster.AddonError); ok && err.IsCritical {
			return err
		}
		log.Warnf(ctx, "Failed to deploy DNS addon execute job for provider coredns: %v", err)
	}
	return nil
}

func doCoreDNSDeploy(ctx context.Context, c *cluster.Cluster) error {
	log.Infof(ctx, "[DNS] Setting up DNS provider %s", c.DNS.Provider)
	CoreDNSConfig := coredns.CoreDNSOptions{
		CoreDNSImage:           c.SystemImages.CoreDNS,
		CoreDNSAutoScalerImage: c.SystemImages.CoreDNSAutoscaler,
		RBACConfig:             c.Authorization.Mode,
		ClusterDomain:          c.ClusterDomain,
		ClusterDNSServer:       c.ClusterDNSServer,
		UpstreamNameservers:    c.DNS.UpstreamNameservers,
		ReverseCIDRs:           c.DNS.ReverseCIDRs,
	}
	coreDNSYaml, err := templates.CompileTemplateFromMap(coredns.CoreDNSTemplate, CoreDNSConfig)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, coreDNSYaml, CoreDNSResourceName, false); err != nil {
		return err
	}
	log.Infof(ctx, "[DNS] DNS provider coredns deployed successfully")
	return nil
}
