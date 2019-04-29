package network

import (
	"context"
	"fmt"

	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/network/calico"
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

	Image            = "Image"
	CNIImage         = "CNIImage"
	NodeImage        = "NodeImage"
	ControllersImage = "ControllersImage"
)

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
