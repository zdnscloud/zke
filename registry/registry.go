package registry

import (
	"context"

	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/registry/resources"
	"github.com/zdnscloud/zke/templates"
)

const (
	AdminServerDeployJobName  = "zke-registry-adminserver-deploy-job"
	ChartMuseumDeployJobName  = "zke-registry-chartmuseum-deploy-job"
	ClairDeployJobName        = "zke-registry-clair-deploy-job"
	CoreDeployJobName         = "zke-registry-core-deploy-job"
	DatabaseDeployJobName     = "zke-registry-database-deploy-job"
	IngressDeployJobName      = "zke-registry-ingress-deploy-job"
	JobserviceDeployJobName   = "zke-registry-jobservice-deploy-job"
	NotaryServerDeployJobName = "zke-registry-notaryserver-deploy-job"
	NotarySignerDeployJobName = "zke-registry-notarysigner-deploy-job"
	PortalDeployJobName       = "zke-registry-portal-deploy-job"
	RedisDeployJobName        = "zke-registry-redis-deploy-job"
	RegistryDeployJobName     = "zke-registry-registry-deploy-job"
)

func DeployRegistry(ctx context.Context, c *cluster.Cluster) error {
	if c.Registry.Isenabled == false {
		log.Infof(ctx, "[Registry] Not enable registry plugin, skip it")
		return nil
	}
	log.Infof(ctx, "[Registry] Setting up Registry Plugin")
	config := map[string]interface{}{
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
	}

	if err := doOneDeploy(ctx, c, config, resources.RedisTemplate, RedisDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, config, resources.DatabaseTemplate, DatabaseDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, config, resources.CoreTemplate, CoreDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, config, resources.RegistryTemplate, RegistryDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, config, resources.NotaryServerTemplate, NotaryServerDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, config, resources.NotarySignerTemplate, NotarySignerDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, config, resources.ChartMuseumTemplate, ChartMuseumDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, config, resources.ClairTemplate, ClairDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, config, resources.JobserviceTemplate, JobserviceDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, config, resources.PortalTemplate, PortalDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, config, resources.AdminServerTemplate, AdminServerDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, config, resources.IngressTemplate, IngressDeployJobName); err != nil {
		return err
	}
	return nil

}

func doOneDeploy(ctx context.Context, c *cluster.Cluster, config map[string]interface{}, resourcesTemplate string, deployJobName string) error {
	configYaml, err := templates.CompileTemplateFromMap(resourcesTemplate, config)
	if err != nil {
		return err
	}

	if err := c.DoAddonDeploy(ctx, configYaml, deployJobName, true); err != nil {
		return err
	}
	return nil
}
