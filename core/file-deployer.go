package core

import (
	"context"
	"fmt"
	"path"

	"github.com/zdnscloud/zke/pkg/docker"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/types"

	"github.com/docker/docker/api/types/container"
)

const (
	ContainerName = "file-deployer"
	ServiceName   = "file-deploy"
	ConfigEnv     = "FILE_DEPLOY"
)

func deployFile(ctx context.Context, uniqueHosts []*hosts.Host, alpineImage string, prsMap map[string]types.PrivateRegistry, fileName, fileContents string) error {
	for _, host := range uniqueHosts {
		log.Infof(ctx, "[%s] Deploying file '%s' to node [%s]", ServiceName, fileName, host.Address)
		if err := doDeployFile(ctx, host, fileName, fileContents, alpineImage, prsMap); err != nil {
			return fmt.Errorf("Failed to deploy file '%s' on node [%s]: %v", host.Address, fileName, err)
		}
	}
	return nil
}

func doDeployFile(ctx context.Context, host *hosts.Host, fileName, fileContents, alpineImage string, prsMap map[string]types.PrivateRegistry) error {
	// remove existing container. Only way it's still here is if previous deployment failed
	if err := docker.DoRemoveContainer(ctx, host.DClient, ContainerName, host.Address); err != nil {
		return err
	}
	containerEnv := []string{ConfigEnv + "=" + fileContents}
	imageCfg := &container.Config{
		Image: alpineImage,
		Cmd: []string{
			"sh",
			"-c",
			fmt.Sprintf("t=$(mktemp); echo -e \"$%s\" > $t && mv $t %s && chmod 644 %s", ConfigEnv, fileName, fileName),
		},
		Env: containerEnv,
	}
	hostCfg := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/etc/kubernetes:z", path.Join(host.PrefixPath, "/etc/kubernetes")),
		},
	}
	if err := docker.DoRunContainer(ctx, host.DClient, imageCfg, hostCfg, ContainerName, host.Address, ServiceName, prsMap); err != nil {
		return err
	}
	if err := docker.DoRemoveContainer(ctx, host.DClient, ContainerName, host.Address); err != nil {
		return err
	}
	log.Debugf(ctx, "[%s] Successfully deployed file '%s' on node [%s]", ServiceName, fileName, host.Address)
	return nil
}
