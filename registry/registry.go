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

type RegistryImage struct {
	HarborAdminserver  string `yaml:"harbor_adminserver" json:"harbor_adminserver"`
	HarborChartmuseum  string `yaml:"harbor_chartmuseum" json:"harbor_chartmuseum"`
	HarborClair        string `yaml:"harbor_clair" json:"harbor_clair"`
	HarborCore         string `yaml:"harbor_core" json:"harbor_core"`
	HarborDatabase     string `yaml:"harbor_database" json:"harbor_database"`
	HarborJobservice   string `yaml:"harbor_jobservice" json:"harbor_jobservice"`
	HarborNotaryServer string `yaml:"harbor_notaryserver" json:"harbor_notaryserver"`
	HarborNotarySigner string `yaml:"harbor_notarysigner" json:"harbor_notarysigner"`
	HarborPortal       string `yaml:"harbor_portal" json:"harbor_portal"`
	HarborRedis        string `yaml:"harbor_redis" json:"harbor_redis"`
	HarborRegistry     string `yaml:"harbor_registry" json:"harbor_registry"`
	HarborRegistryctl  string `yaml:"harbor_registryctl" json:"harbor_registryctl"`
}

var DefaultImage = RegistryImage{
	HarborAdminserver:  "goharbor/harbor-adminserver:v1.7.5",
	HarborChartmuseum:  "goharbor/chartmuseum-photon:v0.8.1-v1.7.5",
	HarborClair:        "goharbor/clair-photon:v2.0.8-v1.7.5",
	HarborCore:         "goharbor/harbor-core:v1.7.5",
	HarborDatabase:     "goharbor/harbor-db:v1.7.5",
	HarborJobservice:   "goharbor/harbor-jobservice:v1.7.5",
	HarborNotaryServer: "goharbor/notary-server-photon:v0.6.1-v1.7.5",
	HarborNotarySigner: "goharbor/notary-signer-photon:v0.6.1-v1.7.5",
	HarborPortal:       "goharbor/harbor-portal:v1.7.5",
	HarborRedis:        "goharbor/redis-photon:v1.7.5",
	HarborRegistry:     "goharbor/registry-photon:v2.6.2-v1.7.5",
	HarborRegistryctl:  "goharbor/harbor-registryctl:v1.7.5",
}

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
		"RedisImage":              DefaultImage.HarborRedis,
		"RedisDiskCapacity":       c.Registry.RedisDiskCapacity,
		"DatabaseImage":           DefaultImage.HarborDatabase,
		"DatabaseDiskCapacity":    c.Registry.DatabaseDiskCapacity,
		"CoreImage":               DefaultImage.HarborCore,
		"RegistryImage":           DefaultImage.HarborRegistry,
		"RegistryctlImage":        DefaultImage.HarborRegistryctl,
		"RegistryDiskCapacity":    c.Registry.RegistryDiskCapacity,
		"NotaryServerImage":       DefaultImage.HarborNotaryServer,
		"NotarySignerImage":       DefaultImage.HarborNotarySigner,
		"ChartmuseumImage":        DefaultImage.HarborChartmuseum,
		"ChartmuseumDiskCapacity": c.Registry.ChartmuseumDiskCapacity,
		"ClairImage":              DefaultImage.HarborClair,
		"JobserviceImage":         DefaultImage.HarborJobservice,
		"JobserviceDiskCapacity":  c.Registry.JobserviceDiskCapacity,
		"PortalImage":             DefaultImage.HarborPortal,
		"AdminserverImage":        DefaultImage.HarborAdminserver,
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
