package lvm

import (
	"context"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/k8s"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/storage/common"
)

func doLVMDDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storage agent lvmd")
	cli, err := k8s.GetK8sClientFromConfig("./kube_config_cluster.yml")
	if err != nil {
		return err
	}
	cfg := map[string]interface{}{
		"RBACConfig":       c.Authorization.Mode,
		"StorageLvmdImage": c.SystemImages.StorageLvmd,
		"LabelKey":         common.StorageHostLabels,
		"LabelValue":       StorageType,
		"StorageNamespace": common.StorageNamespace,
	}
	return k8s.DoDeployFromTemplate(cli, LVMDTemplate, cfg)
}

func doLVMStorageDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storageclass lvm")
	cli, err := k8s.GetK8sClientFromConfig("./kube_config_cluster.yml")
	if err != nil {
		return err
	}
	cfg := map[string]interface{}{
		"RBACConfig":                     c.Authorization.Mode,
		"StorageLvmAttacherImage":        c.SystemImages.StorageLvmAttacher,
		"StorageLvmProvisionerImage":     c.SystemImages.StorageLvmProvisioner,
		"StorageLvmDriverRegistrarImage": c.SystemImages.StorageLvmDriverRegistrar,
		"StorageLvmCSIImage":             c.SystemImages.StorageLvmCSI,
		"LabelKey":                       common.StorageHostLabels,
		"LabelValue":                     StorageType,
		"StorageClassName":               StorageClassName,
		"StorageNamespace":               common.StorageNamespace,
	}
	return k8s.DoDeployFromTemplate(cli, LVMStorageTemplate, cfg)
}
