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

var adminServerConfig = map[string]interface{}{}
var chartMuseumConfig = map[string]interface{}{}
var clairConfig = map[string]interface{}{}
var coreConfig = map[string]interface{}{}
var databaseConfig = map[string]interface{}{}
var ingressConfig = map[string]interface{}{}
var jobserviceConfig = map[string]interface{}{}
var notaryServerConfig = map[string]interface{}{}
var notarySignerConfig = map[string]interface{}{}
var portalConfig = map[string]interface{}{}
var redisConfig = map[string]interface{}{}
var registryConfig = map[string]interface{}{}

func DeployRegistry(ctx context.Context, c *cluster.Cluster) error {
	log.Infof(ctx, "[Registry] Setting up Registry Plugin")
	if err := doOneDeploy(ctx, c, redisConfig, resources.RedisTemplate, RedisDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, databaseConfig, resources.DatabaseTemplate, DatabaseDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, coreConfig, resources.CoreTemplate, CoreDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, registryConfig, resources.RegistryTemplate, RegistryDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, notaryServerConfig, resources.NotaryServerTemplate, NotaryServerDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, notarySignerConfig, resources.NotarySignerTemplate, NotarySignerDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, chartMuseumConfig, resources.ChartMuseumTemplate, ChartMuseumDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, clairConfig, resources.ClairTemplate, ClairDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, jobserviceConfig, resources.JobserviceTemplate, JobserviceDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, portalConfig, resources.PortalTemplate, PortalDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, adminServerConfig, resources.AdminServerTemplate, AdminServerDeployJobName); err != nil {
		return err
	}

	if err := doOneDeploy(ctx, c, ingressConfig, resources.IngressTemplate, IngressDeployJobName); err != nil {
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
