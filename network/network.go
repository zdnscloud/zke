package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/network/calico"
	"github.com/zdnscloud/zke/network/coredns"
	"github.com/zdnscloud/zke/network/flannel"
	"github.com/zdnscloud/zke/network/ingress"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
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

	IngressResourceName = "zke-ingress-plugin"

	Image            = "Image"
	CNIImage         = "CNIImage"
	NodeImage        = "NodeImage"
	ControllersImage = "ControllersImage"

	DeployNamespace = "kube-system"
)

func DeployNetwork(ctx context.Context, c *core.Cluster) error {
	if err := DeployNetworkPlugin(ctx, c); err != nil {
		return err
	}

	if err := deployDNSPlugin(ctx, c); err != nil {
		return err
	}

	if err := deployIngressPlugin(ctx, c); err != nil {
		return err
	}
	return nil
}

func DeployNetworkPlugin(ctx context.Context, c *core.Cluster) error {
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

func doFlannelDeploy(ctx context.Context, c *core.Cluster) error {
	flannelConfig := map[string]interface{}{
		ClusterCIDR:      c.ClusterCIDR,
		Image:            c.SystemImages.Flannel,
		CNIImage:         c.SystemImages.FlannelCNI,
		FlannelInterface: c.Network.Options[FlannelIface],
		FlannelBackend: map[string]interface{}{
			"Type":          c.Network.Options[FlannelBackendType],
			"Directrouting": c.Network.Options[FlannelBackendDirectrouting],
		},
		RBACConfig:        c.Authorization.Mode,
		ClusterVersion:    core.GetTagMajorVersion(c.Version),
		"DeployNamespace": DeployNamespace,
	}
	//pluginYaml, err := templates.GetManifest(flannelConfig, FlannelNetworkPlugin)
	pluginYaml, err := templates.CompileTemplateFromMap(flannel.FlannelTemplate, flannelConfig)
	if err != nil {
		return err
	}
	return c.DoAddonDeploy(ctx, pluginYaml, NetworkPluginResourceName, true)
}

func doCalicoDeploy(ctx context.Context, c *core.Cluster) error {
	clientConfig := pki.GetConfigPath(pki.KubeNodeCertName)
	calicoConfig := map[string]interface{}{
		KubeCfg:           clientConfig,
		ClusterCIDR:       c.ClusterCIDR,
		CNIImage:          c.SystemImages.CalicoCNI,
		NodeImage:         c.SystemImages.CalicoNode,
		Calicoctl:         c.SystemImages.CalicoCtl,
		CloudProvider:     c.Network.Options[CalicoCloudProvider],
		RBACConfig:        c.Authorization.Mode,
		"DeployNamespace": DeployNamespace,
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

func deployDNSPlugin(ctx context.Context, c *core.Cluster) error {
	if err := doCoreDNSDeploy(ctx, c); err != nil {
		if err, ok := err.(*core.AddonError); ok && err.IsCritical {
			return err
		}
		log.Warnf(ctx, "Failed to deploy DNS addon execute job for provider coredns: %v", err)
	}
	return nil
}

func doCoreDNSDeploy(ctx context.Context, c *core.Cluster) error {
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

func deployIngressPlugin(ctx context.Context, c *core.Cluster) error {
	if err := doIngressDeploy(ctx, c); err != nil {
		if err, ok := err.(*core.AddonError); ok && err.IsCritical {
			return err
		}
		log.Warnf(ctx, "Failed to deploy addon execute job [%s]: %v", IngressResourceName, err)
	}
	return nil
}

func doIngressDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[ingress] Setting up %s ingress controller", c.Ingress.Provider)
	ingressConfig := ingress.IngressOptions{
		RBACConfig:     c.Authorization.Mode,
		Options:        c.Ingress.Options,
		NodeSelector:   c.Ingress.NodeSelector,
		ExtraArgs:      c.Ingress.ExtraArgs,
		IngressImage:   c.SystemImages.Ingress,
		IngressBackend: c.SystemImages.IngressBackend,
	}
	// since nginx ingress controller 0.16.0, it can be run as non-root and doesn't require privileged anymore.
	// So we can use securityContext instead of setting privileges via initContainer.
	ingressSplits := strings.SplitN(c.SystemImages.Ingress, ":", 2)
	if len(ingressSplits) == 2 {
		version := strings.Split(ingressSplits[1], "-")[0]
		if version < "0.16.0" {
			ingressConfig.AlpineImage = c.SystemImages.Alpine
		}
	}
	// Currently only deploying nginx ingress controller
	ingressYaml, err := templates.CompileTemplateFromMap(ingress.NginxIngressTemplate, ingressConfig)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, ingressYaml, IngressResourceName, false); err != nil {
		return err
	}
	log.Infof(ctx, "[ingress] ingress controller %s deployed successfully", c.Ingress.Provider)
	return nil
}
