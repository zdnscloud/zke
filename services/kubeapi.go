package services

import (
	"context"

	"github.com/zdnscloud/zke/docker"
	"github.com/zdnscloud/zke/hosts"
	"github.com/zdnscloud/zke/pki"
	"github.com/zdnscloud/zke/types"
)

func runKubeAPI(ctx context.Context, host *hosts.Host, df hosts.DialerFactory, prsMap map[string]types.PrivateRegistry, kubeAPIProcess types.Process, alpineImage string, certMap map[string]pki.CertificatePKI) error {

	imageCfg, hostCfg, healthCheckURL := GetProcessConfig(kubeAPIProcess)
	if err := docker.DoRunContainer(ctx, host.DClient, imageCfg, hostCfg, KubeAPIContainerName, host.Address, ControlRole, prsMap); err != nil {
		return err
	}
	if err := runHealthcheck(ctx, host, KubeAPIContainerName, df, healthCheckURL, certMap); err != nil {
		return err
	}
	return createLogLink(ctx, host, KubeAPIContainerName, ControlRole, alpineImage, prsMap)
}

func removeKubeAPI(ctx context.Context, host *hosts.Host) error {
	return docker.DoRemoveContainer(ctx, host.DClient, KubeAPIContainerName, host.Address)
}

func RestartKubeAPI(ctx context.Context, host *hosts.Host) error {
	return docker.DoRestartContainer(ctx, host.DClient, KubeAPIContainerName, host.Address)
}
