package zcloud

import (
	"context"

	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/pkg/k8s"
	"github.com/zdnscloud/zke/pkg/log"

	"github.com/zdnscloud/gok8s/client"
	clusteragent "github.com/zdnscloud/zke/zcloud/cluster-agent"
	nodeagent "github.com/zdnscloud/zke/zcloud/node-agent"
	zcloudsa "github.com/zdnscloud/zke/zcloud/sa"
)

const (
	RBACConfig               = "RBACConfig"
	Image                    = "Image"
	NodeAgentPort            = "80"
	ClusterAgentResourceName = "cluster-agent"
	SAResourceName           = "sa"
	ClusterAgentJobName      = "zcloud-cluster-agent"
	SAJobName                = "zcloud-sa"

	StorageNFSProvisionerImage = "StorageNFSProvisionerImage"
)

func DeployZcloudManager(ctx context.Context, c *core.Cluster) error {
	k8sClient, err := k8s.GetK8sClientFromConfig(pki.KubeAdminConfigName)
	if err != nil {
		return err
	}
	if err := doSADeploy(ctx, c, k8sClient); err != nil {
		return err
	}
	if err := doClusterAgentDeploy(ctx, c, k8sClient); err != nil {
		return err
	}
	if err := doNodeAgentDeploy(ctx, c, k8sClient); err != nil {
		return err
	}
	return nil
}

func doSADeploy(ctx context.Context, c *core.Cluster, cli client.Client) error {
	log.Infof(ctx, "[zcloud] Setting up ZcloudSADeploy : %s", SAResourceName)
	saconfig := map[string]interface{}{
		RBACConfig: c.Authorization.Mode,
	}
	return k8s.DoCreateFromTemplate(cli, zcloudsa.SATemplate, saconfig)
}

func doClusterAgentDeploy(ctx context.Context, c *core.Cluster, cli client.Client) error {
	log.Infof(ctx, "[zcloud] Setting up ClusterAgentDeploy : %s", ClusterAgentResourceName)
	clusteragentConfig := map[string]interface{}{
		Image: c.Image.ClusterAgent,
	}
	return k8s.DoCreateFromTemplate(cli, clusteragent.ClusterAgentTemplate, clusteragentConfig)
}

func doNodeAgentDeploy(ctx context.Context, c *core.Cluster, cli client.Client) error {
	log.Infof(ctx, "[zcloud] Setting up NodeAgent")
	cfg := map[string]interface{}{
		"Image":         c.Image.NodeAgent,
		"NodeAgentPort": NodeAgentPort,
	}
	return k8s.DoCreateFromTemplate(cli, nodeagent.NodeAgentTemplate, cfg)
}
