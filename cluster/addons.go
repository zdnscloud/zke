package cluster

import (
	"context"
	"fmt"
	"github.com/zdnscloud/zke/addons"
	"github.com/zdnscloud/zke/pkg/k8s"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/templates"
	"strings"
	"time"
)

const (
	UserAddonResourceName          = "zke-user-addon"
	IngressAddonResourceName       = "zke-ingress-controller"
	UserAddonsIncludeResourceName  = "zke-user-includes-addons"
	DNSAddonResourceName           = "zke-coredns-addon"
	IngressAddonJobName            = "zke-ingress-controller-deploy-job"
	MetricsServerAddonJobName      = "zke-metrics-addon-deploy-job"
	MetricsServerAddonResourceName = "zke-metrics-addon"
	NginxIngressAddonAppName       = "ingress-nginx"
	CoreDNSAddonAppName            = "coredns"
)

type ingressOptions struct {
	RBACConfig     string
	Options        map[string]string
	NodeSelector   map[string]string
	ExtraArgs      map[string]string
	AlpineImage    string
	IngressImage   string
	IngressBackend string
}

type MetricsServerOptions struct {
	RBACConfig         string
	Options            map[string]string
	MetricsServerImage string
	Version            string
}
type CoreDNSOptions struct {
	RBACConfig             string
	CoreDNSImage           string
	CoreDNSAutoScalerImage string
	ClusterDomain          string
	ClusterDNSServer       string
	ReverseCIDRs           []string
	UpstreamNameservers    []string
	NodeSelector           map[string]string
}

type AddonError struct {
	err        string
	IsCritical bool
}

func (e *AddonError) Error() string {
	return e.err
}

func getAddonResourceName(addon string) string {
	AddonResourceName := "zke-" + addon + "-addon"
	return AddonResourceName
}

func (c *Cluster) deployK8sAddOns(ctx context.Context) error {
	if err := c.deployDNS(ctx); err != nil {
		if err, ok := err.(*AddonError); ok && err.IsCritical {
			return err
		}
		log.Warnf(ctx, "Failed to deploy DNS addon execute job for provider %s: %v", DNSAddonResourceName, err)

	}

	if err := c.deployMetricServer(ctx); err != nil {
		if err, ok := err.(*AddonError); ok && err.IsCritical {
			return err
		}
		log.Warnf(ctx, "Failed to deploy addon execute job [%s]: %v", MetricsServerAddonResourceName, err)
	}

	if err := c.deployIngress(ctx); err != nil {
		if err, ok := err.(*AddonError); ok && err.IsCritical {
			return err
		}
		log.Warnf(ctx, "Failed to deploy addon execute job [%s]: %v", IngressAddonResourceName, err)

	}
	return nil
}

func (c *Cluster) deployDNS(ctx context.Context) error {
	log.Infof(ctx, "[DNS] Setting up DNS provider %s", c.DNS.Provider)
	CoreDNSConfig := CoreDNSOptions{
		CoreDNSImage:           c.SystemImages.CoreDNS,
		CoreDNSAutoScalerImage: c.SystemImages.CoreDNSAutoscaler,
		RBACConfig:             c.Authorization.Mode,
		ClusterDomain:          c.ClusterDomain,
		ClusterDNSServer:       c.ClusterDNSServer,
		UpstreamNameservers:    c.DNS.UpstreamNameservers,
		ReverseCIDRs:           c.DNS.ReverseCIDRs,
	}
	coreDNSYaml, err := templates.GetManifest(CoreDNSConfig, c.DNS.Provider)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, coreDNSYaml, getAddonResourceName(c.DNS.Provider), false); err != nil {
		return err
	}
	log.Infof(ctx, "[DNS] DNS provider %s deployed successfully", c.DNS.Provider)
	return nil
}

func (c *Cluster) deployMetricServer(ctx context.Context) error {
	log.Infof(ctx, "[addons] Setting up %s", c.Monitoring.MetricsProvider)
	s := strings.Split(c.SystemImages.MetricsServer, ":")
	versionTag := s[len(s)-1]
	MetricsServerConfig := MetricsServerOptions{
		MetricsServerImage: c.SystemImages.MetricsServer,
		RBACConfig:         c.Authorization.Mode,
		Options:            c.Monitoring.MetricsOptions,
		Version:            GetTagMajorVersion(versionTag),
	}
	metricsYaml, err := templates.GetManifest(MetricsServerConfig, c.Monitoring.MetricsProvider)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, metricsYaml, MetricsServerAddonResourceName, false); err != nil {
		return err
	}
	log.Infof(ctx, "[addons] %s deployed successfully", c.Monitoring.MetricsProvider)
	return nil
}

func (c *Cluster) DoAddonDeploy(ctx context.Context, addonYaml, resourceName string, IsCritical bool) error {
	addonUpdated, err := c.StoreAddonConfigMap(ctx, addonYaml, resourceName)
	if err != nil {
		return &AddonError{fmt.Sprintf("Failed to save addon ConfigMap: %v", err), IsCritical}
	}
	log.Infof(ctx, "[addons] Executing deploy job %s", resourceName)
	k8sClient, err := k8s.NewClient(c.LocalKubeConfigPath, c.K8sWrapTransport)
	if err != nil {
		return &AddonError{fmt.Sprintf("%v", err), IsCritical}
	}
	node, err := k8s.GetNode(k8sClient, c.ControlPlaneHosts[0].HostnameOverride)
	if err != nil {
		return &AddonError{fmt.Sprintf("Failed to get Node [%s]: %v", c.ControlPlaneHosts[0].HostnameOverride, err), IsCritical}
	}
	addonJob, err := addons.GetAddonsExecuteJob(resourceName, node.Name, c.Services.KubeAPI.Image)
	if err != nil {
		return &AddonError{fmt.Sprintf("Failed to generate addon execute job: %v", err), IsCritical}
	}

	if err = c.ApplySystemAddonExecuteJob(addonJob, addonUpdated); err != nil {
		return &AddonError{fmt.Sprintf("%v", err), IsCritical}
	}
	return nil
}

func (c *Cluster) StoreAddonConfigMap(ctx context.Context, addonYaml string, addonName string) (bool, error) {
	log.Infof(ctx, "[addons] Saving ConfigMap for addon %s to Kubernetes", addonName)
	updated := false
	kubeClient, err := k8s.NewClient(c.LocalKubeConfigPath, c.K8sWrapTransport)
	if err != nil {
		return updated, err
	}
	timeout := make(chan bool, 1)
	go func() {
		for {
			updated, err = k8s.UpdateConfigMap(kubeClient, []byte(addonYaml), addonName)
			if err != nil {
				time.Sleep(time.Second * 5)
				continue
			}
			log.Infof(ctx, "[addons] Successfully saved ConfigMap for addon %s to Kubernetes", addonName)
			timeout <- true
			break
		}
	}()
	select {
	case <-timeout:
		return updated, nil
	case <-time.After(time.Second * UpdateStateTimeout):
		return updated, fmt.Errorf("[addons] Timeout waiting for kubernetes to be ready")
	}
}

func (c *Cluster) ApplySystemAddonExecuteJob(addonJob string, addonUpdated bool) error {
	if err := k8s.ApplyK8sSystemJob(addonJob, c.LocalKubeConfigPath, c.K8sWrapTransport, c.AddonJobTimeout, addonUpdated); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) deployIngress(ctx context.Context) error {
	log.Infof(ctx, "[ingress] Setting up %s ingress controller", c.Ingress.Provider)
	ingressConfig := ingressOptions{
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
	ingressYaml, err := templates.GetManifest(ingressConfig, c.Ingress.Provider)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, ingressYaml, IngressAddonResourceName, false); err != nil {
		return err
	}
	log.Infof(ctx, "[ingress] ingress controller %s deployed successfully", c.Ingress.Provider)
	return nil
}
