package cluster

import (
	"context"

	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/templates"
)

const (
	ZcloudPredeployResourceName = "zcloud-predeploy"
	ZcloudPredeployJobName      = "zcloud-predeploy-job"
)

func (c *Cluster) deployZcloudPre(ctx context.Context) error {
	log.Infof(ctx, "[ZcloudPre] Starting up ZcloudPredeply")
	if err := c.doZcloudPreDeploy(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) doZcloudPreDeploy(ctx context.Context) error {
	config := map[string]interface{}{
		RBACConfig: c.Authorization.Mode,
	}
	configYaml, err := templates.GetManifest(config, ZcloudPredeployResourceName)
	if err != nil {
		return err
	}
	if err := c.doAddonDeploy(ctx, configYaml, ZcloudPredeployJobName, true); err != nil {
		return err
	}
	return nil
}
