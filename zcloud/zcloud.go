package zcloud

import (
	"context"

	"github.com/zdnscloud/zke/cluster"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
	clusteragent "github.com/zdnscloud/zke/zcloud/cluster-agent"
	zcloudsa "github.com/zdnscloud/zke/zcloud/sa"
)

const (
	RBACConfig = "RBACConfig"
	Image      = "Image"

	ClusterAgentResourceName = "cluster-agent"
	SAResourceName           = "sa"
	ClusterAgentJobName      = "zcloud-cluster-agent"
	SAJobName                = "zcloud-sa"

	StorageNFSProvisionerImage = "StorageNFSProvisionerImage"
)

func DeployZcloudManager(ctx context.Context, c *cluster.Cluster) error {
	if err := doSADeploy(ctx, c); err != nil {
		return err
	}
	if err := doClusterAgentDeploy(ctx, c); err != nil {
		return err
	}
	return nil
}

func doSADeploy(ctx context.Context, c *cluster.Cluster) error {
	log.Infof(ctx, "[zcloud] Setting up ZcloudSADeploy : %s", SAResourceName)
	saconfig := map[string]interface{}{
		RBACConfig: c.Authorization.Mode,
	}
	configYaml, err := templates.CompileTemplateFromMap(zcloudsa.SATemplate, saconfig)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, configYaml, SAJobName, true); err != nil {
		return err
	}
	return nil
}

func doClusterAgentDeploy(ctx context.Context, c *cluster.Cluster) error {
	log.Infof(ctx, "[zcloud] Setting up ClusterAgentDeploy : %s", ClusterAgentResourceName)
	clusteragentConfig := map[string]interface{}{
		Image: c.SystemImages.ClusterAgent,
	}
	clusteragentYaml, err := templates.CompileTemplateFromMap(clusteragent.ClusterAgentTemplate, clusteragentConfig)
	if err != nil {
		return err
	}
	if err := c.DoAddonDeploy(ctx, clusteragentYaml, ClusterAgentJobName, true); err != nil {
		return err
	}
	return nil
}
