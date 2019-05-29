package lvm

import (
	"context"
	"github.com/zdnscloud/zke/core"
	"github.com/zdnscloud/zke/pkg/log"
	"github.com/zdnscloud/zke/pkg/templates"
	"github.com/zdnscloud/zke/storage/common"
)

func doLVMDDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storage agent lvmd")
	cfg := map[string]interface{}{
		"RBACConfig":       c.Authorization.Mode,
		"StorageLvmdImage": c.SystemImages.StorageLvmd,
		"LabelKey":         common.StorageHostLabels,
		"LabelValue":       StorageType,
		"StorageNamespace": common.StorageNamespace,
	}
	yaml, err := templates.CompileTemplateFromMap(LVMDTemplate, cfg)
	if err != nil {
		return err
	}
	return c.DoAddonDeploy(ctx, yaml, "zke-storage-agent-lvmd", true)
}

func doLVMStorageDeploy(ctx context.Context, c *core.Cluster) error {
	log.Infof(ctx, "[storage] Setting up storageclass lvm")
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
	yaml, err := templates.CompileTemplateFromMap(LVMStorageTemplate, cfg)
	if err != nil {
		return err
	}
	return c.DoAddonDeploy(ctx, yaml, "zke-storage-lvm", true)
}
