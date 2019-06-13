package registry

import (
	"context"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/k8s"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/registry/components"

	"github.com/zdnscloud/gok8s/client"
)

const (
	RegistryCertsCN = "harbor"
	DeployNamespace = "zcloud"
)

var componentsTemplates = map[string]string{
	"redis":         components.RedisTemplate,
	"database":      components.DatabaseTemplate,
	"core":          components.CoreTemplate,
	"registry":      components.RegistryTemplate,
	"notary-server": components.NotaryServerTemplate,
	"notary-signer": components.NotarySignerTemplate,
	"chartmuseum":   components.ChartMuseumTemplate,
	"chair":         components.ClairTemplate,
	"jobservice":    components.JobserviceTemplate,
	"portal":        components.PortalTemplate,
	"adminserver":   components.AdminServerTemplate,
	"ingress":       components.IngressTemplate,
}

func DeployRegistry(ctx context.Context, c *core.Cluster) error {
	if c.Registry.Isenabled {
		log.Infof(ctx, "[Registry] Setting up Registry Plugin")

		templateConfig, k8sClient, registryCert, err := prepare(c)
		if err != nil {
			return err
		}

		for component, template := range componentsTemplates {
			err := k8s.DoDeployFromTemplate(k8sClient, template, templateConfig)
			if err != nil {
				log.Infof(ctx, "[Registry] component %s deploy failed", component)
				return err
			}
		}

		if err := deployRegistryCert(ctx, c, registryCert); err != nil {
			return err
		}

		log.Infof(ctx, "[Registry] Successfully deployed Registry Plugin")
		return nil
	}
	return nil
}

func prepare(c *core.Cluster) (map[string]interface{}, client.Client, string, error) {
	IngresscaCert, IngresstlsCert, IngresstlsKey, err := generateRegistryCerts(c, RegistryCertsCN)
	if err != nil {
		return nil, nil, "", err
	}

	templateConfig := map[string]interface{}{
		"RedisImage":              c.SystemImages.HarborRedis,
		"RedisDiskCapacity":       c.Registry.RedisDiskCapacity,
		"DatabaseImage":           c.SystemImages.HarborDatabase,
		"DatabaseDiskCapacity":    c.Registry.DatabaseDiskCapacity,
		"CoreImage":               c.SystemImages.HarborCore,
		"RegistryImage":           c.SystemImages.HarborRegistry,
		"RegistryctlImage":        c.SystemImages.HarborRegistryctl,
		"RegistryDiskCapacity":    c.Registry.RegistryDiskCapacity,
		"NotaryServerImage":       c.SystemImages.HarborNotaryServer,
		"NotarySignerImage":       c.SystemImages.HarborNotarySigner,
		"ChartmuseumImage":        c.SystemImages.HarborChartmuseum,
		"ChartmuseumDiskCapacity": c.Registry.ChartmuseumDiskCapacity,
		"ClairImage":              c.SystemImages.HarborClair,
		"JobserviceImage":         c.SystemImages.HarborJobservice,
		"JobserviceDiskCapacity":  c.Registry.JobserviceDiskCapacity,
		"PortalImage":             c.SystemImages.HarborPortal,
		"AdminserverImage":        c.SystemImages.HarborAdminserver,
		"RegistryIngressURL":      c.Registry.RegistryIngressURL,
		"NotaryIngressURL":        c.Registry.NotaryIngressURL,
		"IngresscaCertBase64":     getB64Cert(IngresscaCert),
		"IngresstlsCertBase64":    getB64Cert(IngresstlsCert),
		"IngresstlsKeyBase64":     getB64Cert(IngresstlsKey),
		"DeployNamespace":         DeployNamespace,
	}
	k8sClient, err := k8s.GetK8sClientFromConfig(c.LocalKubeConfigPath)
	return templateConfig, k8sClient, IngresscaCert, err
}
