package core

import (
	"context"

	"github.com/zdnscloud/zke/core/pki"
	"github.com/zdnscloud/zke/core/services"
	"github.com/zdnscloud/zke/pkg/hosts"
	"github.com/zdnscloud/zke/pkg/k8s"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/types"

	"github.com/zdnscloud/cement/errgroup"
)

func (c *Cluster) ClusterRemove(ctx context.Context) error {
	if err := c.CleanupNodes(ctx); err != nil {
		return err
	}
	c.CleanupFiles(ctx)
	return nil
}

func cleanUpHosts(ctx context.Context, cpHosts, workerHosts, etcdHosts, storageHosts, edgeHosts []*hosts.Host, cleanerImage string, prsMap map[string]types.PrivateRegistry, externalEtcd bool, storageMap map[string][]string) error {
	uniqueHosts := hosts.GetUniqueHostList(cpHosts, workerHosts, etcdHosts, storageHosts, edgeHosts)

	_, err := errgroup.Batch(uniqueHosts, func(h interface{}) (interface{}, error) {
		runHost := h.(*hosts.Host)
		return nil, runHost.CleanUpAll(ctx, cleanerImage, prsMap, externalEtcd, storageMap)
	})
	return err
}

func (c *Cluster) CleanupNodes(ctx context.Context) error {
	externalEtcd := false
	if len(c.Services.Etcd.ExternalURLs) > 0 {
		externalEtcd = true
	}
	// Remove Worker Plane
	if err := services.RemoveWorkerPlane(ctx, c.WorkerHosts, true); err != nil {
		return err
	}
	// Remove Contol Plane
	if err := services.RemoveControlPlane(ctx, c.ControlPlaneHosts, true); err != nil {
		return err
	}

	// Remove Etcd Plane
	if !externalEtcd {
		if err := services.RemoveEtcdPlane(ctx, c.EtcdHosts, true); err != nil {
			return err
		}
	}

	storageMap := c.getStorageConfigMap()

	// Clean up all hosts
	return cleanUpHosts(ctx, c.ControlPlaneHosts, c.WorkerHosts, c.EtcdHosts, c.StorageHosts, c.EdgeHosts, c.SystemImages.Alpine, c.PrivateRegistriesMap, externalEtcd, storageMap)
}

func (c *Cluster) CleanupFiles(ctx context.Context) error {
	pki.RemoveAdminConfig(ctx, c.LocalKubeConfigPath)
	removeStateFile(ctx, c.StateFilePath)
	return nil
}

func (c *Cluster) RemoveOldNodes(ctx context.Context) error {
	kubeClient, err := k8s.NewClient(c.LocalKubeConfigPath, c.K8sWrapTransport)
	if err != nil {
		return err
	}
	nodeList, err := k8s.GetNodeList(kubeClient)
	if err != nil {
		return err
	}
	uniqueHosts := hosts.GetUniqueHostList(c.EtcdHosts, c.ControlPlaneHosts, c.WorkerHosts, c.StorageHosts, c.EdgeHosts)
	for _, node := range nodeList.Items {
		if k8s.IsNodeReady(node) {
			continue
		}
		host := &hosts.Host{}
		host.HostnameOverride = node.Name
		if !hosts.IsNodeInList(host, uniqueHosts) {
			if err := k8s.DeleteNode(kubeClient, node.Name, c.CloudProvider.Name); err != nil {
				log.Warnf(ctx, "Failed to delete old node [%s] from kubernetes")
			}
		}
	}
	return nil
}

func (c *Cluster) getStorageConfigMap() map[string][]string {
	storageConfigMap := map[string][]string{}
	for _, i := range c.Storage.Lvm {
		storageConfigMap[i.Host] = i.Devs
	}
	for _, i := range c.Storage.Nfs {
		storageConfigMap[i.Host] = i.Devs
	}
	for _, i := range c.Storage.Ceph {
		storageConfigMap[i.Host] = i.Devs
	}
	return storageConfigMap
}
